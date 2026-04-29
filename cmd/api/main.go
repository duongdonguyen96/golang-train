package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"

	"golang-train/internal/auth/application/port"
	authuc "golang-train/internal/auth/application/usecase"
	authhttp "golang-train/internal/auth/delivery/http"
	authmw "golang-train/internal/auth/middleware"

	identityport "golang-train/internal/identity/application/port"
	identityuc "golang-train/internal/identity/application/usecase"
	identityhttp "golang-train/internal/identity/delivery/http"
	identityinfra "golang-train/internal/identity/infrastructure/persistence"

	paymentport "golang-train/internal/payment/application/port"
	paymentuc "golang-train/internal/payment/application/usecase"
	paymenthttp "golang-train/internal/payment/delivery/http"
	paymentfileinfra "golang-train/internal/payment/infrastructure/file"
	paymentinfra "golang-train/internal/payment/infrastructure/persistence"
	paymentqueryinfra "golang-train/internal/payment/infrastructure/query"
	paymentstorageinfra "golang-train/internal/payment/infrastructure/storage"

	"golang-train/internal/shared/config"
	"golang-train/internal/shared/database"
	sharedmw "golang-train/internal/shared/middleware"
	sharedhttp "golang-train/internal/shared/delivery/http"

	_ "golang-train/docs"
)

func main() {
	cfg := config.LoadFromEnv()

	resolver := tenant.NewHeaderResolver("X-Tenant-ID")
	factory := database.NewMySQLFactory(cfg.DB)

	// ---- Identity wiring (write port)
	var userRepo identityport.UserRepository = identityinfra.NewGormUserRepository()

	// ---- Payment wiring
	var paymentRepo paymentport.PaymentRepository = paymentinfra.NewGormPaymentRepository()
	var paymentQuery paymentport.PaymentQuery = paymentqueryinfra.NewGormPaymentQuery()
	var exporter paymentport.FileExporter = paymentfileinfra.NewCSVExporter()
	var storage paymentport.Storage = paymentstorageinfra.NewS3LikeStorage(cfg.S3)

	// ---- Auth wiring (query port for credentials)
	var credQuery port.UserCredentialsQuery = identityinfra.NewGormUserCredentialQuery()
	loginUC := authuc.NewLoginUsecase(credQuery, []byte(cfg.JWT.Secret))

	getUserUC := identityuc.NewGetUserUsecase(userRepo)
	createPaymentUC := paymentuc.NewCreatePaymentUsecase(paymentRepo)
	listPaymentsUC := paymentuc.NewListPaymentsUsecase(paymentQuery)
	exportPaymentsUC := paymentuc.NewExportPaymentsUsecase(paymentQuery, exporter, storage)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Multi-tenant DB injection (NO global DB)
	r.Use(sharedmw.TenantDBMiddleware(resolver, factory))

	// Public
	authhttp.RegisterRoutes(r, loginUC)

	sharedhttp.RegisterHello(r)

	// Protected
	jwtMW := authmw.NewJWTMiddleware([]byte(cfg.JWT.Secret))
	api := r.Group("/api")
	api.Use(jwtMW.MustAuth())
	identityhttp.RegisterRoutes(api, getUserUC)
	paymenthttp.RegisterRoutes(api, createPaymentUC, listPaymentsUC, exportPaymentsUC)

	log.Printf("listening on %s", cfg.AppAddr)
	if err := r.Run(cfg.AppAddr); err != nil {
		log.Fatal(err)
	}
}
