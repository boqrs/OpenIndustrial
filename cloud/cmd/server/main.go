package main

import (
	"context"
	"log"

	"github.com/OpenIndustrial/cloud/internal/api"
	"github.com/OpenIndustrial/cloud/internal/config"
	"github.com/OpenIndustrial/cloud/internal/factory"
	"github.com/OpenIndustrial/cloud/internal/device"
	"github.com/OpenIndustrial/cloud/internal/kernel/resource"
	"github.com/OpenIndustrial/cloud/internal/kernel/security"
	"github.com/OpenIndustrial/cloud/internal/kernel/security/provider"

	"github.com/OpenIndustrial/cloud/internal/persistence/postgres"

	"github.com/gin-gonic/gin"
)

// ==================================================================================
// 占位符实现 (Mock Implementations)
// ==================================================================================
// mockMQTT 是 MQTTProvider 依赖的占位符。
type mockMQTT struct{}

func (m *mockMQTT) Endpoint() string { return "ssl://dummy-mqtt.local" }
func (m *mockMQTT) Port() int        { return 8883 }
func (m *mockMQTT) Protocol() string { return "mqtt" }

// mockTxManager 是 TransactionManager 依赖的占位符。
type mockTxManager struct{}

func (m *mockTxManager) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	log.Println("警告: 正在使用 mock Transaction Manager，操作不具备事务性。")
	return fn(ctx)
}

// AuthMiddleware 是认证中间件的占位符。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// ==================================================================================
// 主程序 (Main Application)
// ==================================================================================

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 2. 连接数据库
	db, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 3. 初始化所有仓储 (Repositories)
	resourceRepo := postgres.NewResourceRepository(db)
	attrDefRepo := postgres.NewAttributeDefinitionRepository(db)
	resAttrRepo := postgres.NewResourceAttributeRepository(db)
	resConnRepo := postgres.NewResourceConnectionsRepository(db)
	permissionRepo := postgres.NewPermissionRepository(db)
	credRepo := postgres.NewCredentialRepository(db)
	identityRepo := postgres.NewIdentityRepository(db)
	certRepo := postgres.NewCertificateRepository(db)
	factoryRepo := postgres.NewFactoryRepository(db)
	deviceRepo := postgres.NewDeviceRepository(db)
	dtRepo := postgres.NewDeviceTypeRepository(db)
	// --- Initialize Certificate Authority --
	pkiConfig := provider.ProviderConfig{
		Provider: "",
		AWS: provider.AWSConfig{
			Region:                 cfg.PKI.AWS.Region,
			AccessKey:                  cfg.PKI.AWS.AccessKey,
			SecretKey: cfg.PKI.AWS.SecretKey,
			CAArn:           cfg.PKI.AWS.CAArn,
		},

		Aliyun: provider.AliyunConfig{
			Endpoint:       cfg.PKI.Aliyun.Endpoint,
			AccessKeyID:    cfg.PKI.Aliyun.AccessKeyID,
			AccessKeySecret:    cfg.PKI.Aliyun.AccessKeySecret,
			ParentIdentifier: cfg.PKI.Aliyun.ParentIdentifier,
		},
	}

	pkiFactory, err := provider.NewFactory(pkiConfig)
	if err != nil {
		log.Fatalf("failed to create pki provider factory: %v", err)
	}

		// 2. Create the specific Certificate Authority instance based on config
	ca, err := pkiFactory.Create(provider.Provider(cfg.PKI.Provider))
	if err != nil {
		log.Fatalf("failed to create certificate authority: %v", err)
	}

	adaptedCA := security.NewCertificateAuthorityAdapter(ca)


	// 4. 初始化所有服务 (Services)
	resourceSvc := resource.NewService(resourceRepo, attrDefRepo, resAttrRepo, resConnRepo)
	securitySvc := security.NewService(
		resourceRepo,
		credRepo,
		identityRepo,
		certRepo,
		adaptedCA,
		&mockMQTT{},
		&mockTxManager{},
	)
	factorySvc := factory.NewService(resourceSvc, factoryRepo)
	deviceSvc := device.NewService(resourceSvc, dtRepo, deviceRepo)

	// 5. 设置 HTTP 服务器和处理器 (Handlers)
	router := gin.Default()

	resourceHandler := api.NewResourceHandler(resourceSvc, permissionRepo, AuthMiddleware())
	securityHandler := api.NewSecurityHandler(securitySvc)
	factoryHandler := api.NewFactoryAPI(factorySvc)
	deviceHandler := api.NewDeviceAPI(deviceSvc)

	// 6. 注册路由
	// 【修复】: `resourceHandler.RegisterRoutes` 需要一个 `*gin.RouterGroup`。
	// 我们创建一个 API group (例如 /v1) 并将其传入。
	// `securityHandler` 仍然接收 `*gin.Engine`，因为它没有报错。
	apiV1Group := router.Group("/v1")
	resourceHandler.RegisterRoutes(apiV1Group)
	securityHandler.RegisterSecurityRoutes(apiV1Group)
	factoryHandler.RegisterRouts(apiV1Group)
	deviceHandler.RegisterRoutes(apiV1Group)

	// 7. 启动服务器
	// 【修复】: `cfg.Address` 未定义。使用一个默认的占位符地址来让程序可以运行。
	// TODO: 请将 "0.0.0.0:8080" 替换为您 `config.Config` 结构体中正确的服务器地址字段。
	listenAddress := "0.0.0.0:8080"
	log.Printf("server starting on %s", listenAddress)
	if err := router.Run(listenAddress); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}