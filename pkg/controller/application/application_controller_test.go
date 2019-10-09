package application

import (
	"context"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	api "github.com/bigkevmcd/applications/pkg/apis/app/v1alpha1"
)

var _ reconcile.Reconciler = &ReconcileApplication{}

const (
	namespace = "testing"
	appName   = "test-application"
)

var (
	testConfig = map[string]string{"testing": "value"}
)

func TestCreateUnknownApplication(t *testing.T) {
	r, cl := createApplicationReconciler(t, makeApplication())
	req := makeRequest()

	res, err := r.Reconcile(req)

	fatalIfError(t, "failed to reconcile", err)
	if res.Requeue != false {
		t.Fatalf("res.Requeue got %v, wanted %v", res.Requeue, false)
	}

	cm := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: appName + "-config", Namespace: namespace}, cm)
	if err != nil {
		t.Fatalf("failed to get created config-map: %s", err)
	}
	if !reflect.DeepEqual(cm.Data, testConfig) {
		t.Fatalf("got %#v, wanted %#v", cm.Data, testConfig)
	}
}

func TestCreateUnknownApplicationWithError(t *testing.T) {
	t.Skip()
}

func TestUpdateExistingApplication(t *testing.T) {
	t.Skip()
}

func createApplicationReconciler(t *testing.T, obj ...runtime.Object) (ReconcileApplication, client.Client) {
	scheme := createFakeScheme(t, obj...)
	cl := fake.NewFakeClientWithScheme(scheme, obj...)
	return ReconcileApplication{
		client: cl,
		scheme: scheme,
	}, cl
}

func createFakeScheme(t *testing.T, objs ...runtime.Object) *runtime.Scheme {
	registerObjs := objs
	registerObjs = append(registerObjs, &corev1.ConfigMap{})
	api.SchemeBuilder.Register(registerObjs...)
	scheme, err := api.SchemeBuilder.Build()
	if err != nil {
		t.Fatalf("unable to build scheme: %s", err)
	}
	return scheme
}

func fatalIfError(t *testing.T, msg string, err error) {
	if err != nil {
		t.Fatalf("%s: %s", msg, err)
	}
}

func makeRequest() reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      appName,
			Namespace: namespace,
		},
	}
}

func makeApplication() *api.Application {
	return &api.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: namespace,
		},
		Spec: api.ApplicationSpec{
			Config: map[string]string{"testing": "value"},
		},
	}
}
