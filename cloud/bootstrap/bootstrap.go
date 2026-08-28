package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	//_ "github.com/boqrs/OpenIndustrial/cloud/docs"
	"github.com/boqrs/nexus/database"
	"github.com/boqrs/nexus/email"
	zlog "github.com/boqrs/nexus/log"
	"github.com/boqrs/nexus/redis"
	"github.com/boqrs/nexus/tracing"
	"github.com/boqrs/zeus/ginx"

	"github.com/boqrs/OpenIndustrial/cloud/config"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
	dSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/event"
	fSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/factory"
	idtSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/identity"
	rSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
	secSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security/provider"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/bom"
	plSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/planning"
	woSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	execSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"

	//mSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/material"
	routSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	pSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/product"

	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/execution"
	"github.com/boqrs/OpenIndustrial/cloud/internal/handlers/middleware"
	ph "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/product"
	rh "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/resource"
	sech "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/security"
	woh "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/wokerorder"

	//plh "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/planning"

	dh "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/device"
	fh "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/factory"
	idth "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/identity"

	//phh "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/manufacturing/planning"
	pdh "github.com/boqrs/OpenIndustrial/cloud/internal/handlers/product"
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
	//log.Println("警告: 正在使用 mock Transaction Manager，操作不具备事务性。")
	return fn(ctx)
}


type InfraCloseFunc func() error

// InitInfra 初始化基础设施并注册路由。
func InitInfra(router ginx.ZeroGinRouter) (InfraCloseFunc, error) {
	cfg, cfgMgr, err := config.InitConfigWithManager()
	if err != nil {
		return nil, fmt.Errorf("init config failed: %w", err)
	}

	// 1. 数据库（支持热更新）
	dbProv, err := database.NewProvider(cfg.DBCfg)
	if err != nil {
		return nil, fmt.Errorf("init db provider: %w", err)
	}
	cfgMgr.AddReloader(dbProv)

	// 2. Redis（支持热更新）
	redisProv, err := redis.NewProvider(&cfg.RedisCfg)
	if err != nil {
		return nil, fmt.Errorf("init redis provider: %w", err)
	}
	cfgMgr.AddReloader(redisProv)

	// 3. 日志（支持热更新）
	logProv, err := zlog.NewProviderWithKey("log_cfg", cfg.LogCfg)
	if err != nil {
		return nil, fmt.Errorf("init log provider: %w", err)
	}
	cfgMgr.AddReloader(logProv)

	// 4. 追踪（支持热更新）
	tracingProv, err := tracing.NewProvider(&cfg.Trace)
	if err != nil {
		return nil, fmt.Errorf("init tracing provider: %w", err)
	}
	cfgMgr.AddReloader(tracingProv)

	// 6. email
	emailProv, err := email.NewProvider(cfg.EmailCfg)
	if err != nil {
		return nil, fmt.Errorf("init email provider: %w", err)
	}
	cfgMgr.AddReloader(emailProv)

	// --- Initialize Certificate Authority --
	pkiConfig := provider.ProviderConfig{
		Provider: "",
		AWS: provider.AWSConfig{
			Region:                 cfg.Ca.AWS.Region,
			AccessKey:                  cfg.Ca.AWS.AccessKey,
			SecretKey: cfg.Ca.AWS.SecretKey,
			CAArn:           cfg.Ca.AWS.CAArn,
		},

		Aliyun: provider.AliyunConfig{
			Endpoint:       cfg.Ca.Aliyun.Endpoint,
			AccessKeyID:    cfg.Ca.Aliyun.AccessKeyID,
			AccessKeySecret:    cfg.Ca.Aliyun.AccessKeySecret,
			ParentIdentifier: cfg.Ca.Aliyun.ParentIdentifier,
		},
	}

	pkiFactory, err := provider.NewFactory(pkiConfig)
	if err != nil {
		log.Fatalf("failed to create pki provider factory: %v", err)
	}

		// 2. Create the specific Certificate Authority instance based on config
	ca, err := pkiFactory.Create(provider.Provider(cfg.Ca.Provider))
	if err != nil {
		log.Fatalf("failed to create certificate authority: %v", err)
	}

	adaptedCA := secSrv.NewCertificateAuthorityAdapter(ca)

	router.Use(tracing.GinMiddleware(cfg.Trace.ServiceName))
	// middleware
	// router.Use(tracing.GinMiddleware(cfg.Trace.ServiceName))
	// swagger doc
	router.Handle(http.MethodGet, "/swagger/*any", func(c *gin.Context) ginx.Render {
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
		return nil
	})

	//持久化模块初始化
	resourceRepo := postgres.NewResourceRepository(dbProv)
 	attrDefRepo := postgres.NewAttributeDefinitionRepository(dbProv)
	resAttrRepo := postgres.NewResourceAttributeRepository(dbProv)
	resConnRepo := postgres.NewResourceConnectionsRepository(dbProv)
	permissionRepo := postgres.NewPermissionRepository(dbProv)
	credRepo := postgres.NewCredentialRepository(dbProv)
	identityRepo := postgres.NewIdentityRepository(dbProv)
	certRepo := postgres.NewCertificateRepository(dbProv)
	factoryRepo := postgres.NewFactoryRepository(dbProv)
	deviceRepo := postgres.NewDeviceRepository(dbProv)
	pRepo := postgres.NewProductRepository(dbProv)
	TenantRepo := postgres.NewTenantRepository(dbProv)
	usRepo	:=	postgres.NewUserRepository(dbProv)
	roleRepo :=	 postgres.NewRoleRepository(dbProv)
	groupRepo := postgres.NewGroupRepository(dbProv)
	ePb		:=  event.NewEventPubSub()
	//plRepo :=	postgres.NewProductionPlanRepository(dbProv)
	auth := middleware.NewAuthService(cfg.UserJwtSecret, permissionRepo)
	woRepo := postgres.NewWorkOrderRepository(dbProv)
	plRepo := postgres.NewProductionPlanRepository(dbProv)
	bomRepo	:= postgres.NewBOMRepository(dbProv)
	materialRepo := postgres.NewMaterialRepository(dbProv)
	ufRepo := postgres.NewUnitOfWork(dbProv)
	routingRepo := postgres.NewRoutingRepository(dbProv)
	execRepo :=	 postgres.NewExecutionRepository(dbProv)


	// 业务模块初始化, service层按道理只能使用srv
	reSrv := rSrv.NewService(resourceRepo, attrDefRepo, resAttrRepo, resConnRepo)
	seSrv := secSrv.NewService(resourceRepo, credRepo, identityRepo, certRepo, adaptedCA, &mockMQTT{}, &mockTxManager{})
	pSrv :=	pSrv.NewService(reSrv, pRepo)
	dSrv := dSrv.NewService(deviceRepo, reSrv,pSrv, seSrv)
	fSrv := fSrv.NewService(reSrv, factoryRepo)
	idtSrv :=	idtSrv.NewService(TenantRepo, usRepo, roleRepo, groupRepo, cfg.UserJwtSecret, ePb)
	plSrv := plSrv.NewService(plRepo, pSrv, fSrv)
	//materialSrv := mSrv.NewService(materialRepo)
	routSrv := routSrv.NewService(routingRepo)
	boSrv := bom.NewService(bomRepo, materialRepo, pSrv, ufRepo)
	woSrv := woSrv.NewService(woRepo, plSrv, boSrv, routSrv)
	execSrv := execSrv.NewService(execRepo, woSrv, routSrv)
   // plSrv := planning.NewService(plRepo, pSrv, fSrv)

	// 初始化中间件，用于创建认证等中间件 需要时传入middlewareFactory
	ph.NewHandler(pSrv).RouterRegister(router)
	rh.NewHandler(reSrv, auth).RouterRegister(router)
	sech.NewHandler(seSrv).RouterRegister(router)
	dh.NewHandler(dSrv).RouterRegister(router)
	fh.NewHandler(fSrv).RouterRegister(router)
	idth.NewIdentityHandler(idtSrv, auth).RouterRegister(router)
	pdh.NewHandler(pSrv).RouterRegister(router)
	//plh.NewHandler(plSrv).RouterRegister(router)
	woh.NewHandler(woSrv, auth).RouterRegister(router)
	execution.NewHandler(execSrv, auth).RouterRegister(router)

	return func() error {
		var errs []error
		//if err := tracingProv.Close(); err != nil {
		//	errs = append(errs, fmt.Errorf("shutdown tracing: %w", err))
		//}

		if err := dbProv.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close db: %w", err))
		}

		if err := redisProv.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close redis: %w", err))
		}
		return nil
	}, nil
}