package main

// e.scaleClient.Scales(scaledObject.Namespace).Get(ctx, scaledObject.Status.ScaleTargetGVKR.GroupResource(), scaledObject.Spec.ScaleTargetRef.Name, metav1.GetOptions{})
// 	_, err := e.scaleClient.Scales(scaledObject.Namespace).Update(ctx, scaledObject.Status.ScaleTargetGVKR.GroupResource(), scale, metav1.UpdateOptions{})
import (
	"fmt"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/scale"
	ctrl "sigs.k8s.io/controller-runtime"

	kedautil "github.com/kedacore/keda/v2/pkg/util"
)

var log = ctrl.Log.WithName("scaleclient")

// InitScaleClient initializes scale client and returns k8s version
func InitScaleClient(mgr ctrl.Manager) (scale.ScalesGetter, kedautil.K8sVersion, error) {
	const op = "InitScaleClient"
	kubeVersion := kedautil.K8sVersion{}

	// Initialize discovery client
	discoveryClient, err := createDiscoveryClient(mgr)
	if err != nil {
		return nil, kubeVersion, wrapError(op, "failed to create discovery client", err)
	}

	// Get Kubernetes version
	kubeVersion, err = getKubernetesVersion(discoveryClient)
	if err != nil {
		return nil, kubeVersion, wrapError(op, "failed to get Kubernetes version", err)
	}

	// Create scale client
	scaleClient, err := createScaleClient(mgr, discoveryClient)
	if err != nil {
		return nil, kubeVersion, wrapError(op, "failed to create scale client", err)
	}

	return scaleClient, kubeVersion, nil
}

// createDiscoveryClient creates a new discovery client
func createDiscoveryClient(mgr ctrl.Manager) (*discovery.DiscoveryClient, error) {
	// TODO: Consider adding QPS configuration here if needed
	client, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		log.Error(err, "failed to create discovery client")
		return nil, err
	}
	return client, nil
}

// getKubernetesVersion retrieves and parses the Kubernetes server version
func getKubernetesVersion(client *discovery.DiscoveryClient) (kedautil.K8sVersion, error) {
	version, err := client.ServerVersion()
	if err != nil {
		log.Error(err, "failed to get Kubernetes version")
		return kedautil.K8sVersion{}, err
	}
	return kedautil.NewK8sVersion(version), nil
}

// createScaleClient creates a new scale client
func createScaleClient(mgr ctrl.Manager, discoveryClient *discovery.DiscoveryClient) (scale.ScalesGetter, error) {
	return scale.New(
		discoveryClient.RESTClient(),
		mgr.GetRESTMapper(),
		dynamic.LegacyAPIPathResolverFunc,
		scale.NewDiscoveryScaleKindResolver(discoveryClient),
	), nil
}

// wrapError adds context to an error
func wrapError(op, msg string, err error) error {
	log.Error(err, msg)
	return fmt.Errorf("%s: %s: %w", op, msg, err)
}

func main() {
	
}