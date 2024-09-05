package main

import (
	"log"

	"kub/controller"
	"kub/pkg"
	"kub/repositories"
	"kub/uc"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// CORS middleware'i ekleyin
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Her yerden erişime izin ver
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Kubernetes client
	var kubClient = initKubernetes()

	var podsRepo = repositories.NewPodsRepository(kubClient)
	var podsUC = uc.NewPodsUC(podsRepo)
	var podsHandlers = controller.NewPodsHandler(podsUC)

	var namespaceRepo = repositories.NewNamespaceRepository(kubClient)
	var namespaceUC = uc.NewNamespaceUC(namespaceRepo)
	var namespaceHandlers = controller.NewNamespaceHandler(namespaceUC)

	var podsRoutes = e.Group("/pods")
	podsRoutes.POST("", podsHandlers.Create)

	var namespaceRoutes = e.Group("/namespace")
	namespaceRoutes.GET("", namespaceHandlers.Get)

	e.Logger.Fatal(e.Start(":8080"))
}

func initKubernetes() *kubernetes.Clientset {
	client, err := pkg.NewKubernetesClient()
	if err != nil {
		panic(err.Error())
	}

	return client
}
