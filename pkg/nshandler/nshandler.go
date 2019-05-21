package nshandler

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// WorkshopNamespacer is interface to abstract k8s namespace creation via WorkshopNamespacer CRs
type WorkshopNamespacer interface {
	// Create workshopnamespace, for namespace and name
	Create(string, string) error
	GetKubeconfig(string) ([]byte, error)
}

// NewCreateHandler creates instance of Handler for creating new WorkshopNamespace
func NewCreateHandler(workshopNamespacer WorkshopNamespacer, authCookie string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// Create/Update WorkshopNamespace
		err = workshopNamespacer.Create(namespace, userCookie.Value)
		if err != nil {
			log.Printf("unable create CR: %v", err)
		}
	})
}

// NewKubeconfigGetHandler creates instance of Handler for retrieving kubeconfig
func NewKubeconfigGetHandler(workshopNamespacer WorkshopNamespacer, authCookie string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace, ok := vars["namespace"]
		if !ok {
			http.Error(w, "Namespace name missing", http.StatusBadRequest)
			return
		}

		// Create/Update WorkshopNamespace
		kubeconfig, err := workshopNamespacer.GetKubeconfig(namespace)
		if err != nil {
			http.Error(w, "Kubeconfig not found, try again later", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=config")
		w.Header().Set("Content-Type", "text/plain")

		http.ServeContent(w, r, "config", time.Now(), bytes.NewReader(kubeconfig))
	})
}
