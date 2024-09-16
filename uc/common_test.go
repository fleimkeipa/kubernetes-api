package uc

import (
	"log"

	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
	"k8s.io/client-go/kubernetes"
)

func loadEnv() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	log.Println(".env file loaded successfully")
}

func initTestKubernetes() *kubernetes.Clientset {
	client, err := pkg.NewKubernetesClient()
	if err != nil {
		panic(err.Error())
	}

	return client
}

func initTestDB() *pg.DB {
	return pkg.NewPSQLClient()
}
