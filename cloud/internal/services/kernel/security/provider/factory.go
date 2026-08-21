package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
)

type Provider string

const (
	ProviderAWS    Provider = "aws"
	ProviderAliyun Provider = "aliyun"
)

type Factory struct {
	awsCA    CertificateAuthority
	aliyunCA CertificateAuthority
}

// NewFactory creates a new provider factory.
// It initializes certificate authority providers based on the given configuration.
func NewFactory( pkiConfig ProviderConfig) (*Factory, error) {
	var awsCA CertificateAuthority
	// The awsSDKConfig parameter is ignored as the caller in main.go passes nil.
	// Instead, we load the default AWS config from the environment (a more robust pattern).
	if pkiConfig.AWS.CAArn != "" { // Check if AWS config is intended
		awsInternalCfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to load aws config from environment: %w", err)
		}

		awsCA = NewAWSCA(awsInternalCfg, AWSCAConfig{
			CertificateTemplateARN: pkiConfig.AWS.CAArn,
			//ValidityDays:           pkiConfig.AWS., //TODO：后期增加
		})
	}

	var aliyunCA CertificateAuthority
	// Initialize Aliyun CA if configured.
	if pkiConfig.Aliyun.AccessKeyID != "" { // Check if Aliyun config is intended
		// In a real implementation, we would create an Aliyun client here,
		// likely loading credentials from the environment, similar to AWS.
		// e.g., aliyunClient, err := createAliyunClientFromEnv()
		// For now, as the client implementation is not provided, we pass nil.
		// This will allow the code to compile but will fail at runtime if Aliyun is used.
		aliyunCA = NewAliyunCA(nil, AliyunCAConfig{
			CAInstanceID: pkiConfig.Aliyun.AccessKeyID,
			// RegionID and ValidityDays could also be populated from pkiConfig if they existed.
		})
	}

	return &Factory{awsCA: awsCA, aliyunCA: aliyunCA}, nil
}

func (f *Factory) Create(provider Provider) (CertificateAuthority, error) {

	switch provider {

	case ProviderAWS:

		if f.awsCA == nil {
			return nil, fmt.Errorf(
				"aws certificate authority is not configured",
			)
		}

		return f.awsCA, nil

	case ProviderAliyun:

		if f.aliyunCA == nil {
			return nil, fmt.Errorf(
				"aliyun certificate authority is not configured",
			)
		}

		return f.aliyunCA, nil

	default:

		return nil, fmt.Errorf(
			"unsupported certificate authority provider: %s",
			provider,
		)
	}
}