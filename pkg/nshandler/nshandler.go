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
func New(workshopNamespacer WorkshopNamespacer, authCookie string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("Method %s is not supported", r.Method), http.StatusBadRequest)
			return
		}

		namespace := r.FormValue("namespace")
		if namespace == "" {
			http.Error(w, "Namespace name missing", http.StatusBadRequest)
			return
		}
		userCookie, err := r.Cookie(authCookie)
		if err != nil {
			http.Error(w, "User Cookie missing", http.StatusBadRequest)
			return
		}
		user := userCookie.Value

		// Create/Update WorkshopNamespace
		err = workshopNamespacer.Create(namespace, user)
		if err != nil {
			log.Printf("unable create CR: %v", err)
		}
	})
}
