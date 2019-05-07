package handler

import (
	"context"
	"log"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	workshopnamespacev1alpha1 "github.com/prgcont/workshop-namespace-operator/pkg/apis/operator/v1alpha1"
)

// New creates instance of Handler
func New(c client.Client) Handler {
	return Handler{
		client: c,
	}
}

// Handler implements http handler
type Handler struct {
	client client.Client
}

// CreateWorkshopNamespace is method implementing ServeHttp and creating WorkshopNamespace in k8s cluster
func (h Handler) CreateWorkshopNamespace(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	namespace := r.Form["namespace"][0]
	log.Printf("Creating WorkshopNamespace %s", namespace)

	wn := workshopnamespacev1alpha1.WorkshopNamespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespace,
			Namespace: "default",
		},
	}
	err := h.client.Create(context.TODO(), &wn)
	if err != nil {
		log.Printf("unable create CR: %v", err)
	}

	// h.client.Create(ctx context.Context, obj runtime.Object)
	// w.Header().Set("Content-Disposition", "attachment; filename=config")
	// w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	// http.ServeFile(w, r, kubecfgFile)
	return
}
