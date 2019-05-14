package nshandler

import (
	"fmt"
	"log"
	"net/http"
)

// WorkshopNamespacer is interface to abstract k8s namespace creation via WorkshopNamespacer CRs
type WorkshopNamespacer interface {
	// Create workshopnamespace, for namespace and name
	Create(string, string) error
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

// ServeHTTP is method implementing ServeHttp and creating WorkshopNamespace in k8s cluster
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Method %s is not supported", r.Method), http.StatusBadRequest)
		return
	}

	namespace := r.FormValue("namespace")
	if namespace == "" {
		http.Error(w, "Namespace name missing", http.StatusBadRequest)
		return
	}
	userCookie, err := r.Cookie(authenticatedUserCookie)
	if err != nil {
		http.Error(w, "User Cookie missing", http.StatusBadRequest)
		return
	}
	user := userCookie.Value

	// Create/Update WorkshopNamespace
	err = h.workshopNamespacer.Create(namespace, user)
	if err != nil {
		log.Printf("unable create CR: %v", err)
	}
}
