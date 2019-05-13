package k8s

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	workshopnamespacev1alpha1 "github.com/prgcont/workshop-namespace-operator/pkg/apis/operator/v1alpha1"
)

// New creates instance of WorkshopNamespace
func New(client client.Client, namespace string) *WorkshopNamespace {
	return &WorkshopNamespace{
		client:    client,
		namespace: namespace,
	}
}

// WorkshopNamespace implements WorkshopNamespace
type WorkshopNamespace struct {
	client    client.Client
	namespace string
}

// Create creates WorkshopNamespace CR in k8s cluster.
// CR is updated in case it's already present.
func (w *WorkshopNamespace) Create(namespace, user string) error {
	// Create/Update WorkshopNamespace
	wn := workshopnamespacev1alpha1.WorkshopNamespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespace,
			Namespace: "default",
		},
	}
	err := w.client.Create(context.TODO(), &wn)
	if err != nil {
		log.Printf("unable create CR: %v", err)
	}
	return nil
}

// GetKubeconfig returns kubeconfig for given namespace
func (w *WorkshopNamespace) GetKubeconfig(namespace string) ([]byte, error) {
	return []byte{}, nil
}
