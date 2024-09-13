package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fleimkeipa/kubernetes-api/controller"
	_ "github.com/fleimkeipa/kubernetes-api/docs" // which is the generated folder after swag init
	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/uc"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"k8s.io/client-go/kubernetes"
)

func main() {
	loadEnv()

	serveApplication()
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println(".env file loaded successfully")
}

func serveApplication() {
	var e = echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Add CORS in middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("ALLOW_ORIGIN")},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Kubernetes client
	var kubClient = initKubernetes()

	var dbClient = initDB()

	var eventRepo = repositories.NewEventRepository(dbClient)

	var podsRepo = repositories.NewPodsRepository(kubClient)
	var podsUC = uc.NewPodsUC(podsRepo, eventRepo)
	var podsHandlers = controller.NewPodsHandler(podsUC)

	var namespaceRepo = repositories.NewNamespaceRepository(kubClient)
	var namespaceUC = uc.NewNamespaceUC(namespaceRepo)
	var namespaceHandlers = controller.NewNamespaceHandler(namespaceUC)

	var deploymentRepo = repositories.NewDeploymentInterfaces(kubClient)
	var deploymentUC = uc.NewDeploymentUC(deploymentRepo)
	var deploymentHandlers = controller.NewDeploymentHandler(deploymentUC)

	var podsRoutes = e.Group("/pods")
	podsRoutes.GET("", podsHandlers.List)
	podsRoutes.POST("", podsHandlers.Create)
	podsRoutes.GET("/:id", podsHandlers.GetByNameOrUID)
	podsRoutes.DELETE("/:id", podsHandlers.Delete)
	podsRoutes.PUT("/:id", podsHandlers.Update)

	var namespaceRoutes = e.Group("/namespaces")
	namespaceRoutes.GET("", namespaceHandlers.Get)
	namespaceRoutes.POST("", namespaceHandlers.Create)
	namespaceRoutes.GET("/:id", namespaceHandlers.GetByNameOrUID)
	namespaceRoutes.DELETE("/:id", namespaceHandlers.Delete)
	namespaceRoutes.PUT("/:id", namespaceHandlers.Update)

	var deploymentRoutes = e.Group("/deployments")
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
