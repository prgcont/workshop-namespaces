package handler_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prgcont/workshop-namespaces/pkg/handler"
)

func newFakeWorkshopNamespacer(kubeconfig []byte) fakeWorkshopNamespacer {
	return fakeWorkshopNamespacer{
		Kubeconfig: kubeconfig,
	}
}

type fakeWorkshopNamespacer struct {
	Namespace  string
	Kubeconfig []byte
}

func (f *fakeWorkshopNamespacer) Create(namespace string) error {
	f.Namespace = namespace
	return nil
}

func (f *fakeWorkshopNamespacer) GetKubeconfig(namespace string) ([]byte, error) {
	return f.Kubeconfig, nil
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

	for _, table := range tt {
		t.Run(table.description, func(r *testing.T) {
			runAssert := assert.New(r)

			// Create Test Request
			req, err := http.NewRequest(
				"POST",
				"/namespace",
				strings.NewReader(table.data.Encode()),
			)
			runAssert.NoError(err, "Test Request can't be created")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			// Initialize handler
			wn := newFakeWorkshopNamespacer([]byte{})
			nsHandler := handler.New(&wn)
			handlerFunc := http.HandlerFunc(nsHandler.WorkshopNamespaceHandler)
			rr := httptest.NewRecorder()
			handlerFunc.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			runAssert.Equal(rr.Code, table.returnCode, "handler returned wrong status code")

			// Check the response body is what we expect.
			runAssert.Equal(table.body, rr.Body.String(), "handler returned unexpected body")
		})
	}
}
