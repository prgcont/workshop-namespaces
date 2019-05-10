package handler

import (
	"log"
	"net/http"
)

// WorkshopNamespacer is interface to abstract k8s namespace creation via WorkshopNamespacer CRs
type WorkshopNamespacer interface {
	Create(string) error
	GetKubeconfig(string) ([]byte, error)
}

// New creates instance of Handler
func New(workshopNamespacer WorkshopNamespacer) Handler {
	return Handler{
		workshopNamespacer: workshopNamespacer,
	}
}

// Handler implements http handler
type Handler struct {
	workshopNamespacer WorkshopNamespacer
}

// WorkshopNamespaceHandler is method implementing ServeHttp and creating WorkshopNamespace in k8s cluster
func (h Handler) WorkshopNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	namespace := r.FormValue("namespace")
	if namespace == "" {
		http.Error(w, "Namespace name missing", http.StatusBadRequest)
		return
	}

	// Create/Update WorkshopNamespace
	err := h.workshopNamespacer.Create(namespace)
	if err != nil {
		log.Printf("unable create CR: %v", err)
	}
}
