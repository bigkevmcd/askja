package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mitchellh/go-homedir"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func getProfile(ctx context.Context, c corev1client.CoreV1Interface, namespace, name string) ([]byte, error) {
	b, err := c.
		Services(namespace).
		ProxyGet("http", name, "8000", "/profiles", map[string]string{"q": "test"}).
		DoRaw(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch certificate: %v", err)
	}
	return b, nil
}

func main() {
	kubeConfigPath, err := homedir.Expand("~/.kube/config")
	if err != nil {
		log.Fatalf("failed to expand dir: %s", err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("Failed to get in cluster config: %v", err)
	}

	cl, err := corev1client.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create a client: %s", err)
	}

	b, err := getProfile(context.Background(), cl, "profiles-system", "profiles-controller-controller-manager-metrics-service")
	if err != nil {
		log.Fatalf("Failed to get profile: %s", err)
	}
	log.Printf("KEVIN!!!! %s\n", b)
}
