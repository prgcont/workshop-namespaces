package nshandler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/prgcont/workshop-namespaces/pkg/nshandler"
)

func newFakeWorkshopNamespacer(kubeconfig map[string][]byte) fakeWorkshopNamespacer {
	return fakeWorkshopNamespacer{
		kubeconfig: kubeconfig,
	}
}

type fakeWorkshopNamespacer struct {
	namespace  string
	kubeconfig map[string][]byte
}

func (f *fakeWorkshopNamespacer) Create(namespace, name string) error {
	f.namespace = namespace
	return nil
}

func (f *fakeWorkshopNamespacer) GetKubeconfig(namespace string) ([]byte, error) {
	config, ok := f.kubeconfig[namespace]
	if !ok {
		return []byte{}, fmt.Errorf("Kubeconfig not found")
	}
	return config, nil
}

func TestWorkshopNamespaceHandler(t *testing.T) {
	tt := []struct {
		description string
		data        url.Values
		returnCode  int
		body        string
	}{
		{
			description: "Namespace is created",
			data:        url.Values{"namespace": {"test"}},
			returnCode:  http.StatusOK,
			body:        "",
		},
		{
			description: "Namespace name is missing",
			data:        url.Values{},
			returnCode:  http.StatusBadRequest,
			body:        "Namespace name missing\n",
		},
	}

	cookieKey := "auth.user"
	for _, table := range tt {
		t.Run(table.description, func(r *testing.T) {
			runAssert := assert.New(r)

			// Create Test Request
			req, err := http.NewRequest(
				"POST",
				"/namespaces",
				strings.NewReader(table.data.Encode()),
			)
			runAssert.NoError(err, "Test Request can't be created")

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			userCookie := &http.Cookie{
				Name:    cookieKey,
				Value:   "dummy",
				Path:    "/",
				Expires: time.Now().Add(time.Hour * 2),
			}
			req.AddCookie(userCookie)

			// Initialize handler
			wn := newFakeWorkshopNamespacer(map[string][]byte{})
			nsHandler := nshandler.NewCreateHandler(&wn, cookieKey)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.Handle("/namespaces", nsHandler).Methods("POST")
			router.ServeHTTP(rr, req)

			runAssert.Equal(table.returnCode, rr.Code, "handler returned wrong status code")
			runAssert.Equal(table.body, rr.Body.String(), "handler returned unexpected body")
		})
	}
}

func TestKubeconfigGetHandler(t *testing.T) {
	tt := []struct {
		description string
		requestPath string
		returnCode  int
		body        string
		kubeconfigs map[string][]byte
	}{
		{
			description: "Kubeconfig is downloaded",
			requestPath: "/kubeconfig/test",
			returnCode:  http.StatusOK,
			body:        "TEST",
			kubeconfigs: map[string][]byte{"test": []byte("TEST")},
		},
		{
			description: "Kubeconfig is missing",
			requestPath: "/kubeconfig/test",
			returnCode:  http.StatusNotFound,
			body:        "Kubeconfig not found, try again later\n",
			kubeconfigs: map[string][]byte{},
		},
	}

	cookieKey := "auth.user"
	for _, table := range tt {
		t.Run(table.description, func(r *testing.T) {
			runAssert := assert.New(r)

			// Create Test Request
			req, err := http.NewRequest(
				"GET",
				table.requestPath,
				strings.NewReader(""),
			)
			runAssert.NoError(err, "Test Request can't be created")

			userCookie := &http.Cookie{
				Name:    cookieKey,
				Value:   "dummy",
				Path:    "/",
				Expires: time.Now().Add(time.Hour * 2),
			}
			req.AddCookie(userCookie)

			// Initialize handler
			wn := newFakeWorkshopNamespacer(table.kubeconfigs)
			kubeconfigGetHandler := nshandler.NewKubeconfigGetHandler(&wn, cookieKey)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.Handle("/kubeconfig/{namespace}", kubeconfigGetHandler).Methods("GET")
			router.ServeHTTP(rr, req)

			runAssert.Equal(table.returnCode, rr.Code, "handler returned wrong status code")
			runAssert.Equal(table.body, rr.Body.String(), "handler returned unexpected body")
		})
	}
}
