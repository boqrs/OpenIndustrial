package bootstrap

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/boqrs/nexus/database"
	"github.com/boqrs/nexus/email"
	zlog "github.com/boqrs/nexus/log"
	"github.com/boqrs/nexus/redis"
	"github.com/boqrs/nexus/tracing"
	"github.com/boqrs/zeus/ginx"

	"github.com/boqrs/OpenIndustrial/cloud/config"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/device"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/execution"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/factory"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/identity"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/product"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/resource"
	sh "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/security"
	wh "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/wokerorder"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"

	dSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/event"
	fSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/factory"
	idtSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/identity"

	rSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
	secSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security/provider"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/application"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/bom"
	execSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution/executors"
	plSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/planning"
	routSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	woSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"

	executionresult "github.com/boqrs/OpenIndustrial/cloud/internal/services/executionresult"
	pSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/product"
)

type InfraCloseFunc func() error

// mockMQTT is a temporary MQTT provider implementation.
//
// MQTT connection management is not part of the current bootstrap scope.
// Keep this adapter until the real MQTT provider is connected.
type mockMQTT struct{}

func (m *mockMQTT) Endpoint() string {
	return "ssl://dummy-mqtt.local"
}

func (m *mockMQTT) Port() int {
	return 8883
}

func (m *mockMQTT) Protocol() string {
	return "mqtt"
}

// InitInfra initializes infrastructure, services and HTTP routes.
func InitInfra(router ginx.ZeroGinRouter) (InfraCloseFunc, error) {
	// =========================================================================
	// 1. Configuration
	// =========================================================================

	cfg, cfgMgr, err := config.InitConfigWithManager()
	if err != nil {
		return nil, fmt.Errorf("init config failed: %w", err)
	}

	// =========================================================================
	// 2. Infrastructure Providers
	// =========================================================================

	// -------------------------------------------------------------------------
	// Database
	// -------------------------------------------------------------------------

	dbProv, err := database.NewProvider(cfg.DBCfg)
	if err != nil {
		return nil, fmt.Errorf("init db provider: %w", err)
	}

	cfgMgr.AddReloader(dbProv)

	// -------------------------------------------------------------------------
	// Redis
	// -------------------------------------------------------------------------

	redisProv, err := redis.NewProvider(&cfg.RedisCfg)
	if err != nil {
		return nil, fmt.Errorf("init redis provider: %w", err)
	}

	cfgMgr.AddReloader(redisProv)

	// -------------------------------------------------------------------------
	// Logger
	// -------------------------------------------------------------------------

	logProv, err := zlog.NewProviderWithKey(
		"log_cfg",
		cfg.LogCfg,
	)
	if err != nil {
		return nil, fmt.Errorf("init log provider: %w", err)
	}

	cfgMgr.AddReloader(logProv)

	// -------------------------------------------------------------------------
	// Tracing
	// -------------------------------------------------------------------------

	tracingProv, err := tracing.NewProvider(&cfg.Trace)
	if err != nil {
		return nil, fmt.Errorf("init tracing provider: %w", err)
	}

	cfgMgr.AddReloader(tracingProv)

	// -------------------------------------------------------------------------
	// Email
	// -------------------------------------------------------------------------

	emailProv, err := email.NewProvider(cfg.EmailCfg)
	if err != nil {
		return nil, fmt.Errorf("init email provider: %w", err)
	}

	cfgMgr.AddReloader(emailProv)

	// =========================================================================
	// 3. Unit Of Work
	// =========================================================================
	//
	// The whole application uses this UnitOfWork implementation.
	//
	// Do NOT create another transaction manager implementation here.
	//
	// Cross-domain business transactions are handled by the application layer.
	// =========================================================================

	uow := postgres.NewUnitOfWork(dbProv)

	// =========================================================================
	// 4. Certificate Authority
	// =========================================================================

	pkiConfig := provider.ProviderConfig{
		Provider: "",

		AWS: provider.AWSConfig{
			Region:    cfg.Ca.AWS.Region,
			AccessKey: cfg.Ca.AWS.AccessKey,
			SecretKey: cfg.Ca.AWS.SecretKey,
			CAArn:     cfg.Ca.AWS.CAArn,
		},

		Aliyun: provider.AliyunConfig{
			Endpoint:         cfg.Ca.Aliyun.Endpoint,
			AccessKeyID:      cfg.Ca.Aliyun.AccessKeyID,
			AccessKeySecret:  cfg.Ca.Aliyun.AccessKeySecret,
			ParentIdentifier: cfg.Ca.Aliyun.ParentIdentifier,
		},
	}

	pkiFactory, err := provider.NewFactory(pkiConfig)
	if err != nil {
		return nil, fmt.Errorf(
			"create certificate authority factory: %w",
			err,
		)
	}

	ca, err := pkiFactory.Create(
		provider.Provider(cfg.Ca.Provider),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create certificate authority: %w",
			err,
		)
	}

	adaptedCA := secSrv.NewCertificateAuthorityAdapter(ca)

	// =========================================================================
	// 5. HTTP Middleware / Swagger
	// =========================================================================

	router.Use(
		tracing.GinMiddleware(
			cfg.Trace.ServiceName,
		),
	)

	router.Handle(
		http.MethodGet,
		"/swagger/*any",
		func(c *gin.Context) ginx.Render {
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
			return nil
		},
	)

	// =========================================================================
	// 6. Persistence Repositories
	// =========================================================================

	// -------------------------------------------------------------------------
	// Kernel / Resource
	// -------------------------------------------------------------------------

	resourceRepo := postgres.NewResourceRepository(dbProv)
	attrDefRepo := postgres.NewAttributeDefinitionRepository(dbProv)
	resAttrRepo := postgres.NewResourceAttributeRepository(dbProv)
	resConnRepo := postgres.NewResourceConnectionsRepository(dbProv)

	// -------------------------------------------------------------------------
	// Security
	// -------------------------------------------------------------------------

	permissionRepo := postgres.NewPermissionRepository(dbProv)
	credentialRepo := postgres.NewCredentialRepository(dbProv)
	identityRepo := postgres.NewIdentityRepository(dbProv)
	certificateRepo := postgres.NewCertificateRepository(dbProv)

	// -------------------------------------------------------------------------
	// Factory / Device / Product
	// -------------------------------------------------------------------------

	factoryRepo := postgres.NewFactoryRepository(dbProv)
	deviceRepo := postgres.NewDeviceRepository(dbProv)
	productRepo := postgres.NewProductRepository(dbProv)

	// -------------------------------------------------------------------------
	// Identity
	// -------------------------------------------------------------------------

	tenantRepo := postgres.NewTenantRepository(dbProv)
	userRepo := postgres.NewUserRepository(dbProv)
	roleRepo := postgres.NewRoleRepository(dbProv)
	groupRepo := postgres.NewGroupRepository(dbProv)

	// -------------------------------------------------------------------------
	// Manufacturing
	// -------------------------------------------------------------------------

	productionPlanRepo := postgres.NewProductionPlanRepository(dbProv)
	workOrderRepo := postgres.NewWorkOrderRepository(dbProv)
	bomRepo := postgres.NewBOMRepository(dbProv)
	materialRepo := postgres.NewMaterialRepository(dbProv)
	routingRepo := postgres.NewRoutingRepository(dbProv)
	executionRepo := postgres.NewExecutionRepository(dbProv)

	// ExecutionResult repository currently exposes NewRepository.
	executionResultRepo := postgres.NewRepository(dbProv)

	// =========================================================================
	// 7. Kernel Services
	// =========================================================================

	// -------------------------------------------------------------------------
	// Resource
	// -------------------------------------------------------------------------

	resourceService := rSrv.NewService(
		resourceRepo,
		attrDefRepo,
		resAttrRepo,
		resConnRepo,
	)

	// -------------------------------------------------------------------------
	// Security
	//
	// IMPORTANT:
	// The last dependency is now UnitOfWork.
	// mockTxManager has been completely removed.
	// -------------------------------------------------------------------------

	securityService := secSrv.NewService(
		resourceRepo,
		credentialRepo,
		identityRepo,
		certificateRepo,
		adaptedCA,
		&mockMQTT{},
		uow,
	)

	// =========================================================================
	// 8. Product / Device / Factory
	// =========================================================================

	productService := pSrv.NewService(
		resourceService,
		productRepo,
	)

	deviceService := dSrv.NewService(
		deviceRepo,
		resourceService,
		productService,
		//securityService,
	)

	factoryService := fSrv.NewService(
		resourceService,
		factoryRepo,
	)

	// =========================================================================
	// 9. Identity
	// =========================================================================

	eventPubSub := event.NewEventPubSub()

	authService := middleware.NewAuthService(
		cfg.UserJwtSecret,
		permissionRepo,
	)

	identityService := idtSrv.NewService(
		tenantRepo,
		userRepo,
		roleRepo,
		groupRepo,
		cfg.UserJwtSecret,
		eventPubSub,
	)

	// =========================================================================
	// 10. Manufacturing - Planning
	// =========================================================================

	planningService := plSrv.NewService(
		productionPlanRepo,
		productService,
		factoryService,
	)

	// =========================================================================
	// 11. Manufacturing - BOM
	// =========================================================================

	bomService := bom.NewService(
		bomRepo,
		materialRepo,
		productService,
		uow,
	)

	// =========================================================================
	// 12. Manufacturing - Routing
	// =========================================================================

	routingService := routSrv.NewService(
		routingRepo,
	)

	// =========================================================================
	// 13. Manufacturing - WorkOrder
	// =========================================================================

	workOrderService := woSrv.NewService(
		workOrderRepo,
		planningService,
		bomService,
		routingService,
	)

	// =========================================================================
	// 14. Manufacturing - Operation Executors
	// =========================================================================
	//
	// Standard executors:
	//
	//   SN_GENERATE
	//   CERTIFICATE_ISSUE
	//
	// SerialNumberGenerator is intentionally nil for now because the project
	// does not yet have a concrete serial-number generator implementation.
	//
	// This means SN_GENERATE will correctly report
	// "serial number generator is not configured" until that dependency
	// is implemented.
	//
	// Do NOT replace this with a fake generator.
	// =========================================================================

	var serialNumberGenerator executors.SerialNumberGenerator

	executorRegistry := executors.BuildOperationExecutorRegistry(
		serialNumberGenerator,
		adaptedCA,
	)

	// =========================================================================
	// 15. Manufacturing - Execution
	// =========================================================================

	executionService := execSrv.NewService(
		executionRepo,
		workOrderService,
		routingService,
		executorRegistry,
		uow,
	)

	// =========================================================================
	// 16. Manufacturing - ExecutionResult
	// =========================================================================

	executionResultService := executionresult.NewService(
		executionResultRepo,
		executionService,
	)

	// =========================================================================
	// 17. Manufacturing Application Service
	// =========================================================================
	//
	// This is the cross-domain orchestration layer.
	//
	// Example:
	//
	// CreateProductionExecution:
	//   WorkOrder
	//       ↓
	//   Resource
	//       ↓
	//   Execution
	//
	// ConfirmExecutionResult:
	//   ExecutionResult
	//       ↓
	//   Execution
	//       ↓
	//   ExecutionOperation identity
	//       ↓
	//   Device
	//       ↓
	//   WorkOrder.CompletedQuantity
	//       ↓
	//   ExecutionResult.Confirmed
	//
	// All cross-domain writes use the same UnitOfWork.
	// =========================================================================

	manufacturingApplicationService := application.NewService(
		uow,
		workOrderRepo,
		routingRepo,
		executionService,
		executionResultRepo,
		executionRepo,
		deviceService,
		resourceService,
	)

	// =========================================================================
	// 18. HTTP Handlers
	// =========================================================================

	// -------------------------------------------------------------------------
	// Product
	// -------------------------------------------------------------------------

	product.NewHandler(
		productService,
	).RouterRegister(router)

	// -------------------------------------------------------------------------
	// Resource
	// -------------------------------------------------------------------------

	resource.NewHandler(
		resourceService,
		authService,
	).RouterRegister(router)

	// -------------------------------------------------------------------------
	// Security
	// -------------------------------------------------------------------------

	sh.NewHandler(
		securityService,
	).RouterRegister(router)

	// -------------------------------------------------------------------------
	// Device
	// -------------------------------------------------------------------------

	device.NewHandler(
		deviceService,
	).RouterRegister(router)

	// -------------------------------------------------------------------------
	// Factory
	// -------------------------------------------------------------------------

	factory.NewHandler(
		factoryService,
	).RouterRegister(router)

	// -------------------------------------------------------------------------
	// Identity
	// -------------------------------------------------------------------------

	identity.NewIdentityHandler(
		identityService,
		authService,
	).RouterRegister(router)

	// -------------------------------------------------------------------------
	// WorkOrder
	// -------------------------------------------------------------------------

	wh.NewHandler(
		workOrderService,
		authService,
	).RouterRegister(router)

	// -------------------------------------------------------------------------
	// Execution
	// -------------------------------------------------------------------------
	//
	// Existing execution handler currently consumes execution.Service directly.
	// Cross-domain operations such as StartProductionExecution /
	// ConfirmExecutionResult are exposed through the manufacturing application
	// service when their corresponding handlers are introduced.
	// -------------------------------------------------------------------------

	execution.NewHandler(
		executionService,
		authService,
	).RouterRegister(router)

	// =========================================================================
	// Keep the application service alive as part of the composition root.
	// =========================================================================
	//
	// It is intentionally not registered to an HTTP handler here yet because
	// the current handler layer does not expose the complete manufacturing
	// application API.
	//
	// Once the manufacturing application handler is introduced, this is the
	// service that should be injected there.
	// =========================================================================

	_ = executionResultService
	_ = manufacturingApplicationService

	// =========================================================================
	// 19. Shutdown
	// =========================================================================

	return func() error {
		var errs []error

		if err := dbProv.Close(); err != nil {
			errs = append(
				errs,
				fmt.Errorf("close db: %w", err),
			)
		}

		if err := redisProv.Close(); err != nil {
			errs = append(
				errs,
				fmt.Errorf("close redis: %w", err),
			)
		}

		if len(errs) > 0 {
			return fmt.Errorf(
				"shutdown infrastructure: %v",
				errs,
			)
		}

		return nil
	}, nil
}
