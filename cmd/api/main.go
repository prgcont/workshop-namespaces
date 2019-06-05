package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/prgcont/workshop-namespace-operator/pkg/apis"
	wnclient "github.com/prgcont/workshop-namespace-operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/prgcont/workshop-namespaces/pkg/k8s"
	"github.com/prgcont/workshop-namespaces/pkg/nshandler"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	authCookieName := "auth.user"
	secretsClient, wnClientset, err := getClient("default")
	if err != nil {
		log.Printf("unable create client: %v", err)
		os.Exit(1)
	}

	// Handlers
	workshopNamespace := k8s.New(secretsClient, wnClientset, "default")
	nsCreateHandler := nshandler.NewCreateHandler(workshopNamespace, authCookieName)
	kubeconfigGetHandler := nshandler.NewKubeconfigGetHandler(workshopNamespace, authCookieName)

	// Routes
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./static"))).Methods("GET")
	r.Handle("/namespaces", nsCreateHandler).Methods("POST")
	r.Handle("/kubeconfig/{namespace}", kubeconfigGetHandler).Methods("GET")

	// Middlewares
	withCookie := nshandler.NewCookieMiddleware(authCookieName, 10, time.Hour*10)
	r.Use(withCookie)

	// Server start
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func getClient(namespace string) (corev1.SecretInterface, wnclient.WorkshopNamespaceInterface, error) {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err := apis.AddToScheme(scheme.Scheme); err != nil {
		log.Printf("unable add APIs to scheme: %v", err)
		return nil, nil, err
	}

	// Core resources client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	secretsClient := clientset.Core().Secrets(namespace)

	// WorkshopNamespace client
	wnClientset, err := wnclient.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return secretsClient, wnClientset.WorkshopNamespaces(namespace), nil
}

// func getClient() (*kubernetes.Clientset, error) {
// 	var config *rest.Config
// 	var err error

// 	if kubeconfig == "" {
// 		log.Printf("using in-cluster configuration")
// 		config, err = rest.InClusterConfig()
// 	} else {
// 		log.Printf("using configuration from '%s'", kubeconfig)
// 		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
// 	}
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	if err := apis.AddToScheme(scheme.Scheme); err != nil {
// 		log.Printf("unable add APIs to scheme: %v", err)
// 		return nil, err
// 	}
// 	workshopnamespacev1alpha1.AddToScheme(scheme.Scheme)
// 	return client.New(cfg, client.Options{})
// }
