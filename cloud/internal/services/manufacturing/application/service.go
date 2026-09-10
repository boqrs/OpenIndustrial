package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/executionresult"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/google/uuid"
)

var (
	ErrWorkOrderNotReleasable    = errors.New("work order is not in a releasable state for execution")
	ErrRoutingNotActive          = errors.New("routing is not active")
	ErrRoutingHasNoOperations    = errors.New("routing has no operations")
	ErrWorkOrderQuantityExceeded = errors.New("the planned quantity for the work order has already been met or exceeded")
	ErrExecutionResultNotFound   = errors.New("execution result not found")
	ErrExecutionNotCompleted     = errors.New("execution is not completed")
	ErrExecutionAlreadyFailed    = errors.New("execution has failed")
	ErrExecutionResultInvalid    = errors.New("execution result is invalid")
)

type service struct {
	uow              postgres.UnitOfWork
	workOrders       workorder.Repository
	routings         routing.Repository
	executions       execution.Service
	executionResults executionresult.Repository
	devices          device.Service
}

// NewService creates a new manufacturing application service.
func NewService(uow postgres.UnitOfWork, workOrders workorder.Repository, routings routing.Repository, executions execution.Service, executionResults executionresult.Repository, devices device.Service) Service {
	return &service{
		uow:              uow,
		workOrders:       workOrders,
		routings:         routings,
		executions:       executions,
		executionResults: executionResults,
		devices:          devices,
	}
}

func (s *service) CreateProductionExecution(
	ctx context.Context,
	workOrderID uint,
	deviceID *uint,
) (*execution.ExecutionResponse, error) {

	req := &execution.CreateExecutionRequest{
		WorkOrderID: workOrderID,
	}

	return s.executions.CreateExecution(ctx, req)
}

// ConfirmExecutionResult confirms the final production result.
// // // This is the manufacturing completion boundary:
//
//	// // ExecutionResult // ↓
//
// // validate Execution // ↓
//
//	// finalize qualified product identities // ↓
//	// update WorkOrder.CompletedQuantity // ↓
//	// confirm ExecutionResult //
//
// // All operations are executed inside one UnitOfWork transaction.
func (s *service) ConfirmExecutionResult(ctx context.Context, executionResultID uint) error {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return errors.New("tenant ID not found in context")
	}
	return s.uow.Execute(ctx, func(txCtx context.Context) error {
		// ------------------------------------------------------------
		// // 1. Load ExecutionResult // ------------------------------------------------------------
		result, err := s.executionResults.GetByIDForUpdateTx(txCtx, tenantID, executionResultID)
		if err != nil {
			return fmt.Errorf("get execution result: %w", err)
		}

		if result == nil {
			return ErrExecutionResultNotFound
		}
		// Idempotent confirmation.
		if result.Status == model.ExecutionResultStatusConfirmed {
			return nil
		}
		if result.Status != model.ExecutionResultStatusDraft {
			return fmt.Errorf("%w: status=%s", ErrExecutionResultInvalid, result.Status)
		}
		// ------------------------------------------------------------ // 2. Validate quantities //
		// ------------------------------------------------------------
		if result.ProducedQuantity < 0 || result.QualifiedQuantity < 0 || result.RejectedQuantity < 0 {
			return fmt.Errorf("%w: negative quantity", ErrExecutionResultInvalid)
		}

		if result.QualifiedQuantity+result.RejectedQuantity != result.ProducedQuantity {
			return fmt.Errorf("%w: produced=%d qualified=%d rejected=%d", ErrExecutionResultInvalid, result.ProducedQuantity, result.QualifiedQuantity, result.RejectedQuantity)
		} // ------------------------------------------------------------ // 3. Load Execution // ------------------------------------------------------------

		exec, err := s.executions.GetExecution(txCtx, result.ExecutionID)
		if err != nil {
			return fmt.Errorf("get execution: %w", err)
		}
		if exec == nil {
			return errors.New("execution not found")
		}

		if exec.Status != model.ProductionExecutionStatusCompleted {
			if exec.Status == model.ProductionExecutionStatusFailed {
				return ErrExecutionAlreadyFailed
			}
			return ErrExecutionNotCompleted
		}
		if exec.WorkOrderID != result.WorkOrderID {
			return fmt.Errorf("%w: execution work order mismatch", ErrExecutionResultInvalid)
		} // ------------------------------------------------------------ // 4. Load WorkOrder // ------------------------------------------------------------
		workOrder, err := s.workOrders.GetByID(txCtx, tenantID, result.WorkOrderID)
		if err != nil {
			return fmt.Errorf("get work order: %w", err)
		}

		if workOrder == nil {
			return errors.New("work order not found")
		}

		if workOrder.Status == model.WorkOrderStatusCancelled {
			return errors.New("cannot confirm result for cancelled work order")
		} // ------------------------------------------------------------ // 5. Validate cumulative quantity // ------------------------------------------------------------
		newCompletedQuantity := workOrder.CompletedQuantity + result.QualifiedQuantity
		if newCompletedQuantity > workOrder.PlannedQuantity {
			return ErrWorkOrderQuantityExceeded
		}
		// ------------------------------------------------------------ // 6. Finalize qualified product identities // ------------------------------------------------------------ // // IMPORTANT: // // We intentionally do NOT create devices based only on // QualifiedQuantity. // // A quantity tells us "97 products passed". // It does not tell us which 97 physical products they are. // // The actual item-level identity comes from ExecutionOperation // results. // // This extraction will be implemented as the next step: // // ExecutionOperation.Result // ↓ // item_key // ↓ // serial_number // ↓ // Device // // For now we explicitly reject a non-zero qualified quantity // without item-level identity data instead of silently creating // incorrect Devices. //
		if result.QualifiedQuantity > 0 {
			operations, err := s.executions.ListOperations(txCtx, exec.ID)
			if err != nil {
				return fmt.Errorf("list execution operations: %w", err)
			}
			items, err := validateQualifiedItems(
				operations,
				result.QualifiedQuantity,
			)

			if err != nil {
				return err
			}

			deviceRequests := make([]*device.CreateDeviceFromExecutionResultRequest, 0, len(items))
			for _, item := range items {
				req := &device.CreateDeviceFromExecutionResultRequest{
					ProductID:         exec.ProductID,
					WorkOrderID:       exec.WorkOrderID,
					ExecutionID:       exec.ID,
					ExecutionResultID: result.ID,
					SerialNumber:      item.SerialNumber,
					HardwareID:        item.HardwareID,
				}
				deviceRequests = append(deviceRequests, req)
			}

			if _, err := s.devices.CreateFromExecutionResultBatchTx(
				txCtx,
				deviceRequests,
			); err != nil {
				return fmt.Errorf(
					"create devices: %w",
					err,
				)
			}
		}
		// ------------------------------------------------------------ // 7. Update WorkOrder // ------------------------------------------------------------
		workOrder.CompletedQuantity = newCompletedQuantity
		if newCompletedQuantity >= workOrder.PlannedQuantity {
			now := time.Now()
			workOrder.Status = model.WorkOrderStatusCompleted
			workOrder.CompletedAt = &now
		}
		if err := s.workOrders.UpdateTx(txCtx, workOrder); err != nil {
			return fmt.Errorf("update work order: %w", err)
		}
		// ------------------------------------------------------------ // 8. Confirm ExecutionResult // ------------------------------------------------------------
		now := time.Now()
		result.Status = model.ExecutionResultStatusConfirmed
		result.ConfirmedAt = &now
		if err := s.executionResults.UpdateTx(txCtx, result); err != nil {
			return fmt.Errorf("confirm execution result: %w", err)
		}
		return nil
	})
}

type qualifiedItem struct {
	ItemKey      string
	SerialNumber string
	HardwareID   string
}

func validateQualifiedItems(
	operations []*execution.OperationResponse,
	qualifiedQuantity int64,
) ([]qualifiedItem, error) {

	if qualifiedQuantity == 0 {
		return nil, nil
	}

	if len(operations) == 0 {
		return nil, fmt.Errorf(
			"%w: execution has no operations",
			ErrExecutionResultInvalid,
		)
	}

	items, err := extractQualifiedItems(operations)
	if err != nil {
		return nil, err
	}

	if int64(len(items)) != qualifiedQuantity {
		return nil, fmt.Errorf(
			"%w: qualified quantity=%d but identity items=%d",
			ErrExecutionResultInvalid,
			qualifiedQuantity,
			len(items),
		)
	}

	return items, nil
}

func extractQualifiedItems(
	operations []*execution.OperationResponse,
) ([]qualifiedItem, error) {

	itemsByKey := make(map[string]*productionItem)

	for _, operation := range operations {
		if operation == nil {
			continue
		}

		// 没有结果的 Operation 不参与生产单元提取。
		if len(operation.Result) == 0 {
			continue
		}

		raw, err := json.Marshal(operation.Result)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: marshal operation %d result: %v",
				ErrExecutionResultInvalid,
				operation.ID,
				err,
			)
		}

		var envelope operationResultEnvelope

		if err := json.Unmarshal(raw, &envelope); err != nil {
			return nil, fmt.Errorf(
				"%w: invalid result of operation %d: %v",
				ErrExecutionResultInvalid,
				operation.ID,
				err,
			)
		}

		for _, resultItem := range envelope.Items {

			if resultItem.ItemKey == "" {
				return nil, fmt.Errorf(
					"%w: operation %d contains empty item_key",
					ErrExecutionResultInvalid,
					operation.ID,
				)
			}

			if resultItem.Data == nil {
				return nil, fmt.Errorf(
					"%w: operation %d item %s has no data",
					ErrExecutionResultInvalid,
					operation.ID,
					resultItem.ItemKey,
				)
			}

			item, exists := itemsByKey[resultItem.ItemKey]

			if !exists {
				item = &productionItem{
					ItemKey: resultItem.ItemKey,
					Data:    make(map[string]any),
				}

				itemsByKey[resultItem.ItemKey] = item
			}

			for key, value := range resultItem.Data {

				oldValue, exists := item.Data[key]

				if exists {
					if !reflect.DeepEqual(oldValue, value) {
						return nil, fmt.Errorf(
							"%w: conflicting value for item %s field %s",
							ErrExecutionResultInvalid,
							resultItem.ItemKey,
							key,
						)
					}

					continue
				}

				item.Data[key] = value
			}
		}
	}

	result := make([]qualifiedItem, 0, len(itemsByKey))

	for _, item := range itemsByKey {

		serialNumber, ok := stringValue(
			item.Data["serial_number"],
		)

		if !ok || serialNumber == "" {
			return nil, fmt.Errorf(
				"%w: item %s has no serial_number",
				ErrExecutionResultInvalid,
				item.ItemKey,
			)
		}

		hardwareID, _ := stringValue(
			item.Data["hardware_id"],
		)

		result = append(result, qualifiedItem{
			ItemKey:      item.ItemKey,
			SerialNumber: serialNumber,
			HardwareID:   hardwareID,
		})
	}

	// map 遍历顺序不稳定。
	// 排序以后，Device 创建顺序稳定，也方便测试。
	sort.Slice(result, func(i, j int) bool {
		return result[i].ItemKey < result[j].ItemKey
	})

	return result, nil
}

func stringValue(value any) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case nil:
		return "", false
	default:
		return "", false
	}
}
