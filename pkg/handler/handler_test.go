package handler_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	workshopnamespacev1alpha1 "github.com/prgcont/workshop-namespace-operator/pkg/apis/operator/v1alpha1"
	"github.com/prgcont/workshop-namespaces/pkg/handler"
)

func TestHealthCheckHandler(t *testing.T) {
	err := workshopnamespacev1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err, "workshopnamespace scheme can't be added to Scheme")

	tt := []struct {
		description string
		data        url.Values
		statusCode  int
		body        string
	}{
		{
			description: "Namespace is created",
			data:        url.Values{"namespace": {"test"}},
			statusCode:  http.StatusOK,
			body:        "",
		},
	}

	for _, table := range tt {
		t.Run(table.description, func(t *testing.T) {
			assert := assert.New(t)

			// Create Test Request
			req, err := http.NewRequest(
				"POST",
				"/namespace",
				strings.NewReader(table.data.Encode()),
			)
			assert.NoError(err, "Test Request can't be created")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			// Initialize handler
			c := fakeclient.NewFakeClient()
			nsHandler := handler.New(c)
			handlerFunc := http.HandlerFunc(nsHandler.WorkshopNamespaceHandler)
			rr := httptest.NewRecorder()
			handlerFunc.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			assert.Equal(rr.Code, table.statusCode, "handler returned wrong status code")

			// Check the response body is what we expect.
			assert.Equal(rr.Body.String(), table.body, "handler returned unexpected body")
		})
	}
}
