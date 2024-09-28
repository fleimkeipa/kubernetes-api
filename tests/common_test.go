package tests

import (
	"context"
	"fmt"
	"log"

	"github.com/fleimkeipa/kubernetes-api/pkg"

	"github.com/go-pg/pg"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	test_db     *pg.DB
	terminateDB = func() {}
)

func initTestKubernetes() *kubernetes.Clientset {
	client, err := pkg.NewKubernetesClient()
	if err != nil {
		log.Fatalf("Failed to init kubernetes client: %v", err)
	}

	createTestNamespace(client)

	return client
}

func deleteTestNamespace(client *kubernetes.Clientset) {
	err := client.CoreV1().Namespaces().Delete(context.TODO(), "test", v1.DeleteOptions{})
	if err != nil {
		log.Fatalf("failed to delete test namespace: %v", err)
	}
}

func createTestNamespace(client *kubernetes.Clientset) {
	namespace := corev1.Namespace{
		TypeMeta: v1.TypeMeta{
			Kind:       "namespace",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec:   corev1.NamespaceSpec{},
		Status: corev1.NamespaceStatus{},
	}
	_, err := client.CoreV1().Namespaces().Create(context.TODO(), &namespace, v1.CreateOptions{})
	existErr := "namespaces \"test\" already exists"
	if err != nil {
		if err.Error() == existErr {
			return
		}
		log.Fatalf("failed to create test namespace: %v", err)
	}
}

func addTempData(data interface{}) error {
	_, err := test_db.Model(data).Insert()
	if err != nil {
		return err
	}

	return nil
}

func clearTable(tableName string) error {
	query := fmt.Sprintf("TRUNCATE %s; DELETE FROM %s", tableName, tableName)
	_, err := test_db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
