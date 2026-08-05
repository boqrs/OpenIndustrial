// drivers/bacnet/codec.go
package bacnet

import (
	"fmt"
)

// decodeValue 将 gobacnet 返回的数据转换成 Runtime 的 Go 类型。
func decodeValue(value any, dataType DataType) (any, error) {

	switch v := value.(type) {

	case bool:
		if dataType != DataTypeBoolean {
			return nil, fmt.Errorf("expect %s but got Boolean", dataType)
		}
		return v, nil

	case uint32: // 对应 UnsignedInteger 和 Enumerated
		switch dataType {
		case DataTypeUint:
			return v, nil
		case DataTypeInt: // 允许 uint32 转换为 int64
			return int64(v), nil
		case DataTypeEnum:
			return v, nil
		default:
			return nil, fmt.Errorf("expect %s but got UnsignedInteger/Enumerated", dataType)
		}

	case int32: // 对应 SignedInteger
		switch dataType {
		case DataTypeInt:
			return int64(v), nil
		default:
			return nil, fmt.Errorf("expect %s but got SignedInteger", dataType)
		}

	case float32: // 对应 Real
		switch dataType {
		case DataTypeFloat:
			return float64(v), nil // 将 float32 提升为 float64
		default:
			return nil, fmt.Errorf("expect %s but got Real", dataType)
		}

	case float64: // 对应 Double
		switch dataType {
		case DataTypeFloat:
			return v, nil
		default:
			return nil, fmt.Errorf("expect %s but got Double", dataType)
		}

	case string: // 对应 CharacterString
		switch dataType {
		case DataTypeString:
			return v, nil
		default:
			return nil, fmt.Errorf("expect %s but got CharacterString", dataType)
		}

	default:
		return nil, fmt.Errorf("unsupported BACnet value %T", value)
	}
}

// encodeValue 将 Runtime 的 Go 类型转换成 gobacnet 类型。
// 由于 gobacnet 的 encoding 包期望 Go 原生类型，此函数主要执行类型验证并按原样返回该值。
func encodeValue(value any, dataType DataType) (any, error) {

	switch dataType {

	case DataTypeBoolean:
		_, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("expect bool for DataTypeBoolean")
		}
		return value, nil

	case DataTypeInt:
		switch value.(type) {
		case int, int32, int64:
			return value, nil
		default:
			return nil, fmt.Errorf("expect integer for DataTypeInt")
		}

	case DataTypeUint:
		switch value.(type) {
		case uint, uint32, uint64:
			return value, nil
		default:
			return nil, fmt.Errorf("expect unsigned integer for DataTypeUint")
		}

	case DataTypeFloat:
		switch value.(type) {
		case float32, float64:
			return value, nil
		default:
			return nil, fmt.Errorf("expect float for DataTypeFloat")
		}

	case DataTypeString:
		_, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expect string for DataTypeString")
		}
		return value, nil

	case DataTypeEnum:
		switch value.(type) {
		case uint32, int: // 允许 int 用于枚举，gobacnet 会将其转换为 uint32
			return value, nil
		default:
			return nil, fmt.Errorf("expect enum (uint32 or int) for DataTypeEnum")
		}

	default:
		return nil, fmt.Errorf("unsupported datatype %s", dataType)
	}
}