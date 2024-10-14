package pkg

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// Determine if we are on local or cluster
	if stage := viper.GetString("stage"); stage == "prod" {
		config, err = getConfigOnCluster()
		if err != nil {
			return nil, fmt.Errorf("failed to get config for cluster stage: %w", err)
		}
	} else {
		config, err = getConfigOnLocal()
		if err != nil {
			return nil, fmt.Errorf("failed to get config for dev stage: %w", err)
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getConfigOnLocal() (*rest.Config, error) {
	kubeconfigPath := ""

	// Use the KUBECONFIG environment variable if it is set
	if envKubeconfig := os.Getenv("KUBECONFIG"); envKubeconfig != "" {
		kubeconfigPath = envKubeconfig
	}

	// Use the default kubeconfig path if not set
	if kubeconfigPath == "" {
		if os.Geteuid() == 0 {
			kubeconfigPath = "/root/.kube/config"
		} else {
			homeDir := os.Getenv("HOME")
			if homeDir != "" {
				kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
			}
		}
	}

	// Let the user override the path with the --kubeconfig flag
	kubeconfig := flag.String("kubeconfig", kubeconfigPath, "absolute path to the kubeconfig file")
	flag.Parse()

	// Build the client configuration
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func getConfigOnCluster() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("Error creating in-cluster config: %v", err)
	}

	return config, nil
}
