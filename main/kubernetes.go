package main

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
)

var (
	clientset *kubernetes.Clientset
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func kubernetesSetup() error {
	configEnv := strings.Split(os.Getenv("KUBECONFIG"), ":")

	var kubeconfig string
	if len(configEnv) > 0 && configEnv[0] != "" {
		kubeconfig = configEnv[0]
	} else {
		kubeconfig = filepath.Join(homeDir(), ".kube", "config")
	}

	var overrides clientcmd.ConfigOverrides
	if len(configEnv) > 1 && configEnv[1] != "" {
		overrides.CurrentContext = configEnv[1]
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&overrides,
	).ClientConfig()

	if err != nil {
		return err
	}

	clientset, err = kubernetes.NewForConfig(config)
	return err
}

func getNamespaces() (*v1.NamespaceList, error) {
	return clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
}
