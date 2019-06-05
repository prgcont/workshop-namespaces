package k8s_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"

	workshopnamespacev1alpha1 "github.com/prgcont/workshop-namespace-operator/pkg/apis/operator/v1alpha1"
	wnfakeclient "github.com/prgcont/workshop-namespace-operator/pkg/client/clientset/versioned/typed/operator/v1alpha1/fake"

	"github.com/prgcont/workshop-namespaces/pkg/k8s"
)

func TestCreate(t *testing.T) {
	createWN := func(name, namespace, user string) workshopnamespacev1alpha1.WorkshopNamespace {
		return workshopnamespacev1alpha1.WorkshopNamespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: workshopnamespacev1alpha1.WorkshopNamespaceSpec{
				Owner: user,
			},
		}
	}

	err := workshopnamespacev1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err, "workshopnamespace scheme can't be added to Scheme")

	tt := []struct {
		description string
		namespace   string
		name        string
		owner       string
		wnExpected  workshopnamespacev1alpha1.WorkshopNamespace
	}{
		{
			description: "Create WorkshopNamespace object",
			namespace:   "default",
			name:        "test",
			owner:       "user",
			wnExpected:  createWN("test", "default", "user"),
		},
	}

	for _, table := range tt {
		t.Run(table.description, func(r *testing.T) {
			runAssert := assert.New(r)

			fakeClientset := fake.NewSimpleClientset()
			secretsFakeClient := fakeClientset.Core().Secrets(table.namespace)
			fake := wnfakeclient.FakeOperatorV1alpha1{&fakeClientset.Fake}
			wnFakeClient := fake.WorkshopNamespaces(table.namespace)

			workshopNamespace := k8s.New(secretsFakeClient, wnFakeClient, table.namespace)
			workshopNamespace.Create(table.name, table.owner)

			wnOut, err := wnFakeClient.Get(table.name, metav1.GetOptions{})
			runAssert.NoError(err, "Object not created")
			runAssert.Equal(&table.wnExpected, wnOut, "Unexpected object")
		})
	}
}

func TestGetKubeconfig(t *testing.T) {
	err := workshopnamespacev1alpha1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err, "workshopnamespace scheme can't be added to Scheme")

	tt := []struct {
		description         string
		namespace           string
		name                string
		expectedErrorString string
		secret              *v1.Secret
		kubeconfigBytes     []byte
	}{
		{
			description:         "Retrieve kubeconfig",
			namespace:           "default",
			name:                "test",
			expectedErrorString: "",
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kubeconfig-test",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"config": []byte("Test DATA"),
				},
			},
			kubeconfigBytes: []byte("Test DATA"),
		},
		{
			description:         "Can't Retrieve Kubeconfig",
			namespace:           "default",
			name:                "test",
			expectedErrorString: "secrets \"kubeconfig-test\" not found",
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bad-name",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"config": []byte("Test DATA"),
				},
			},
			kubeconfigBytes: []byte{},
		},
		{
			description:         "Secrets with missing config data",
			namespace:           "default",
			name:                "test",
			expectedErrorString: "kubeconfig for namespace test is missing in Secret",
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kubeconfig-test",
					Namespace: "default",
				},
			},
			kubeconfigBytes: []byte{},
		},
	}

	for _, table := range tt {
		t.Run(table.description, func(r *testing.T) {
			runAssert := assert.New(r)

			fakeClientset := fake.NewSimpleClientset()
			secretsFakeClient := fakeClientset.Core().Secrets(table.namespace)
			fake := wnfakeclient.FakeOperatorV1alpha1{&fakeClientset.Fake}
			wnFakeClient := fake.WorkshopNamespaces(table.namespace)

			secretsFakeClient.Create(table.secret)

			workshopNamespace := k8s.New(secretsFakeClient, wnFakeClient, table.namespace)
			config, err := workshopNamespace.GetKubeconfig(table.name)

			if table.expectedErrorString != "" {
				runAssert.EqualError(err, table.expectedErrorString, "")
			}

			runAssert.Equal(table.kubeconfigBytes, config, "Got Unexpected Kubeconfig")
		})
	}
}
