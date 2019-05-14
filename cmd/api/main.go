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
	"github.com/prgcont/workshop-namespaces/pkg/k8s"
	"github.com/prgcont/workshop-namespaces/pkg/nshandler"
)

func main() {
	authCookieName := "auth.user"
	wnClient, err := getClient()
	if err != nil {
		log.Printf("unable create client: %v", err)
		os.Exit(1)
	}

	// Handlers
	workshopNamespace := k8s.New(wnClient, "default")
	nsHandler := nshandler.New(workshopNamespace, authCookieName)

	// Middlewares
	withCookie := nshandler.NewCookieMiddleware(authCookieName, 10, time.Hour*10)

	// Routes
	http.Handle("/", withCookie(http.FileServer(http.Dir("./static"))))
	http.Handle("/namespaces", withCookie(nsHandler))

	// Server start
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
