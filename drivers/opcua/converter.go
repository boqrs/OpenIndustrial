package opcua

import "github.com/gopcua/opcua/ua"

func ToQuality(statusCode ua.StatusCode) Quality {
	switch statusCode {
	case ua.StatusGood:
		return QualityGood
	case ua.StatusUncertain:
		return QualityUncertain
	default:
		return QualityBad
	}
}