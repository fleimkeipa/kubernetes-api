package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fleimkeipa/kubernetes-api/config"
	"github.com/fleimkeipa/kubernetes-api/controller"
	_ "github.com/fleimkeipa/kubernetes-api/docs" // which is the generated folder after swag init
	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/uc"
	"github.com/fleimkeipa/kubernetes-api/util"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

func main() {
	loadEnv()

	serveApplication()
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	log.Println(".env file loaded successfully")
}

func serveApplication() {
	var e = echo.New()
	// e.HideBanner = true
	// e.HidePort = true

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	var sugar = logger.Sugar()
	defer sugar.Sync()

	e.Use(pkg.ZapLogger(logger))
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Add CORS in middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("ALLOW_ORIGIN")},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Kubernetes client
	var kubClient = initKubernetes()

	var dbClient = initDB()

	var eventRepo = repositories.NewEventRepository(dbClient)

	var podsRepo = repositories.NewPodsRepository(kubClient)
	var podsUC = uc.NewPodsUC(podsRepo, eventRepo)
	var podsHandlers = controller.NewPodsHandler(podsUC, sugar)

	var namespaceRepo = repositories.NewNamespaceRepository(kubClient)
	var namespaceUC = uc.NewNamespaceUC(namespaceRepo)
	var namespaceHandlers = controller.NewNamespaceHandler(namespaceUC)

	var deploymentRepo = repositories.NewDeploymentInterfaces(kubClient)
	var deploymentUC = uc.NewDeploymentUC(deploymentRepo)
	var deploymentHandlers = controller.NewDeploymentHandler(deploymentUC)

	var userRepo = repositories.NewUserRepository(dbClient)
	var userUC = uc.NewUserUC(userRepo)
	var userHandlers = controller.NewUserHandlers(userUC)

	config.GoogleConfig()

	var authRoutes = e.Group("/auth")
	authRoutes.POST("/login", userHandlers.Login)

	var googleAuthHandler = controller.NewGoogleAuthHandler(userUC)

	var oauthRoutes = authRoutes.Group("")
	oauthRoutes.GET("/google_login", googleAuthHandler.GoogleLogin)
	oauthRoutes.GET("/google_callback", googleAuthHandler.GoogleCallback)

	var restrictedRoutes = e.Group("")
	restrictedRoutes.Use(util.JWTAuth)
	restrictedRoutes.Use(util.JWTAuthViewer)

	var userRoutes = restrictedRoutes.Group("/users")
	userRoutes.POST("", userHandlers.CreateUser)
	userRoutes.PUT("/:id", userHandlers.UpdateUser)

	var podsRoutes = restrictedRoutes.Group("/pods")
	podsRoutes.GET("", podsHandlers.List)
	podsRoutes.POST("", podsHandlers.Create)
	podsRoutes.GET("/:id", podsHandlers.GetByNameOrUID)
	podsRoutes.DELETE("/:id", podsHandlers.Delete)
	podsRoutes.PUT("/:id", podsHandlers.Update)

	var namespaceRoutes = restrictedRoutes.Group("/namespaces")
	namespaceRoutes.GET("", namespaceHandlers.Get)
	namespaceRoutes.POST("", namespaceHandlers.Create)
	namespaceRoutes.GET("/:id", namespaceHandlers.GetByNameOrUID)
	namespaceRoutes.DELETE("/:id", namespaceHandlers.Delete)
	namespaceRoutes.PUT("/:id", namespaceHandlers.Update)

	var deploymentRoutes = restrictedRoutes.Group("/deployments")
	deploymentRoutes.GET("", deploymentHandlers.List)
	deploymentRoutes.POST("", deploymentHandlers.Create)
	deploymentRoutes.GET("/:id", deploymentHandlers.GetByNameOrUID)
	deploymentRoutes.DELETE("/:id", deploymentHandlers.Delete)
	deploymentRoutes.PUT("/:id", deploymentHandlers.Update)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("API_PORT"))))
}

func initKubernetes() *kubernetes.Clientset {
	client, err := pkg.NewKubernetesClient()
	if err != nil {
		panic(err.Error())
	}

	return client
}

func initDB() *pg.DB {
	return pkg.NewPSQLClient()
}
