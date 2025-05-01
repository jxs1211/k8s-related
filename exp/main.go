package main

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

func execInPod(ctx context.Context, config *rest.Config, clientset *kubernetes.Clientset, podName, namespace, containerName string, command []string) error {
	// clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{}) // check if pod exists
	// Create the exec request
	req := clientset.CoreV1().RESTClient().Post().
		Namespace(namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	fmt.Println(req.URL())
	// equivalent to the output of k exec -it php-apache-678865dd57-fpds7 -c php-apache -v 7 -- ls -al
	// Create the executor
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to create SPDY executor: %v", err)
	}

	// Connect to the pod and set up the streams
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return fmt.Errorf("failed to stream: %v", err)
	}
	return nil
}

func main() {
	// Load Kubernetes config
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(fmt.Errorf("failed to build config: %v", err))
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create clientset: %v", err))
	}

	// Example usage
	podName := "php-apache-678865dd57-fpds7"
	namespace := "default"
	containerName := "php-apache"
	command := []string{"ls", "-al"}
	ctx := context.Background()
	err = execInPod(ctx, config, clientset, podName, namespace, containerName, command)
	if err != nil {
		panic(fmt.Errorf("exec failed: %v", err))
	}
}
