package k8s_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	workshopnamespacev1alpha1 "github.com/prgcont/workshop-namespace-operator/pkg/apis/operator/v1alpha1"

	"github.com/prgcont/workshop-namespaces/pkg/k8s"
)

func TestCreate(t *testing.T) {
	createWN := func(name, namespace string) workshopnamespacev1alpha1.WorkshopNamespace {
		return workshopnamespacev1alpha1.WorkshopNamespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			}}
	}

	err := workshopnamespacev1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err, "workshopnamespace scheme can't be added to Scheme")

	tt := []struct {
		description string
		namespace   string
		name        string
		wnExpected  workshopnamespacev1alpha1.WorkshopNamespace
	}{
		{
			description: "Create WorkshopNamespace object",
			namespace:   "default",
			name:        "test",
			wnExpected:  createWN("test", "default"),
		},
	}

	for _, table := range tt {
		t.Run(table.description, func(r *testing.T) {
			runAssert := assert.New(r)

			c := fakeclient.NewFakeClient()
			workshopNamespace := k8s.New(c, table.namespace)
			workshopNamespace.Create(table.name)

			wnOut := workshopnamespacev1alpha1.WorkshopNamespace{}
			err = c.Get(context.Background(), types.NamespacedName{Name: table.name, Namespace: table.namespace}, &wnOut)
			runAssert.NoError(err, "Object not created")
			runAssert.Equal(table.wnExpected, wnOut, "Unexpected object")
		})
	}
}
