package pkg

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
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

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
