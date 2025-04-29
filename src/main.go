package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var filePath string = "/etc/hosts"
var nameSpace string = "default"
var secretName string

func createk8scli() *kubernetes.Clientset {
	// Load in-cluster config if running in a pod, otherwise use kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		kubeconfig := os.Getenv("KUBECONFIG")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(fmt.Sprintf("Error loading kubeconfig or in-cluster config: %v\n", err))
		}
	}
	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("Error creating Kubernetes client: %v\n", err))
	}
	return clientset
}

func main() {
	clientset := createk8scli()

	fp := os.Getenv("FILEPATH")
	if fp != "" {
		filePath = fp
	}
	sn := os.Getenv("SECRETNAME")
	if sn != "" {
		secretName = sn
	} else {
		secretName = strings.ReplaceAll(filePath[1:], "/", "-")
	}
	ns := os.Getenv("NAMESPACE")
	if ns != "" {
		nameSpace = ns
	}

	// fmt.Println(filePath, secretName, nameSpace)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error reading file: %v\n", err))
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: nameSpace,
		},
		Data: map[string][]byte{
			"file-content": []byte(fileContent),
		},
	}

	_, err = clientset.CoreV1().Secrets(nameSpace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("Secret '%s/%s' already exists\n", secretName, nameSpace)
		return
	}

	_, err = clientset.CoreV1().Secrets(nameSpace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(fmt.Sprintf("Error creating secret: %v\n", err))
	}

	fmt.Printf("Secret '%s/%s' created successfully\n", secretName, nameSpace)
}
