package k8s

import (
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	workshopnamespacev1alpha1 "github.com/prgcont/workshop-namespace-operator/pkg/apis/operator/v1alpha1"

	wnclient "github.com/prgcont/workshop-namespace-operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
)

// New creates instance of WorkshopNamespace
func New(secretsClient corev1.SecretInterface, wnclient wnclient.WorkshopNamespaceInterface, namespace string) *WorkshopNamespace {
	return &WorkshopNamespace{
		secretsClient: secretsClient,
		wnclient:      wnclient,
		namespace:     namespace,
	}
}

// WorkshopNamespace implements WorkshopNamespace
type WorkshopNamespace struct {
	secretsClient corev1.SecretInterface
	wnclient      wnclient.WorkshopNamespaceInterface
	namespace     string
}

// Create creates WorkshopNamespace CR in k8s cluster.
// CR is updated in case it's already present.
func (w *WorkshopNamespace) Create(namespace, user string) error {
	// Create/Update WorkshopNamespace
	wn := workshopnamespacev1alpha1.WorkshopNamespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
		Spec: workshopnamespacev1alpha1.WorkshopNamespaceSpec{
			Owner: user,
		},
	}

	_, err := w.wnclient.Create(&wn)
	if err != nil {
		log.Printf("unable create CR: %v", err)
		return err
	}

	return nil
}

// GetKubeconfig returns kubeconfig for given namespace
func (w *WorkshopNamespace) GetKubeconfig(kubeconfigNamespace string) ([]byte, error) {
	// Find Secret for given Namespace
	secret, err := w.secretsClient.Get(fmt.Sprint("kubeconfig-", kubeconfigNamespace), metav1.GetOptions{})
	if err != nil {
		log.Printf("unable retrieve Kubeconfig: %v", err)
		return []byte{}, err
	}

	config, ok := secret.Data["config"]
	if !ok {
		log.Printf("unable retrieve Kubeconfig: %v", err)
		return []byte{}, fmt.Errorf("kubeconfig for namespace %s is missing in Secret", kubeconfigNamespace)
	}

	return config, nil
}
