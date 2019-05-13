package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/prgcont/workshop-namespace-operator/pkg/apis"
	workshopnamespacev1alpha1 "github.com/prgcont/workshop-namespace-operator/pkg/apis/operator/v1alpha1"
	"github.com/prgcont/workshop-namespaces/pkg/handler"
	"github.com/prgcont/workshop-namespaces/pkg/k8s"
)

// TODO: Handle errors: https://blog.questionable.services/article/http-handler-error-handling-revisited/
// TODO: Login with cookies
// TODO: Create and return
// TODO: Watch for Secret to be created
func main() {
	wnClient, err := getClient()
	if err != nil {
		log.Printf("unable create client: %v", err)
		os.Exit(1)
	}

	wn := k8s.New(wnClient, "default")
	h := handler.New(wn)
	cookieAuthHandler := handler.NewCookieAuth(time.Hour * 2)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/namespaces", cookieAuthHandler.CookieAuth(h))
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func getClient() (client.Client, error) {
	// Get config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Printf("unable get config: %v", err)
		os.Exit(1)
	}

	if err := apis.AddToScheme(scheme.Scheme); err != nil {
		log.Printf("unable add APIs to scheme: %v", err)
		return nil, err
	}
	workshopnamespacev1alpha1.AddToScheme(scheme.Scheme)
	return client.New(cfg, client.Options{})
}

// cookieIdentifier Simple cookie identifier to set cookie to identify user, no auth at the moment
func cookieIdentifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Responding to Request: %+v\n", r)
		next.ServeHTTP(w, r)
		return
	})
}
