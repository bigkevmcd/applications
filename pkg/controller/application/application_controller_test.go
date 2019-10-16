package application

import (
	"context"
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	api "github.com/bigkevmcd/applications/pkg/apis/app/v1alpha1"
)

var _ reconcile.Reconciler = &ReconcileApplication{}

const (
	testNamespace = "testing"
	testAppName   = "test-application"
	testReplicas  = 5
)

func TestCreateUnknownApplicationConfiguration(t *testing.T) {
	r, cl := createApplicationReconciler(t, makeTestApplication())
	req := makeRequest()

	res, err := r.Reconcile(req)

	fatalIfError(t, "failed to reconcile", err)
	if res.Requeue != false {
		t.Fatalf("res.Requeue got %v, wanted %v", res.Requeue, false)
	}

	app := &api.Application{}
	err = cl.Get(context.TODO(), ns(testAppName, testNamespace), app)
	if err != nil {
		t.Fatalf("failed to get application: %s", err)
	}
	assertConfigMapHasData(t, testAppName, testNamespace, cl, testEnvironment)
}

func TestCreateUnknownApplicationDeployment(t *testing.T) {
	r, cl := createApplicationReconciler(t, makeTestApplication())
	req := makeRequest()

	res, err := r.Reconcile(req)

	fatalIfError(t, "failed to reconcile", err)
	if res.Requeue != false {
		t.Fatalf("res.Requeue got %v, wanted %v", res.Requeue, false)
	}

	dp := &appsv1.Deployment{}
	err = cl.Get(context.TODO(), ns(testAppName, testNamespace), dp)
	if err != nil {
		t.Fatalf("failed to get created deployment: %s", err)
	}
	assertDeploymentConfiguration(t, testAppName, testNamespace, cl, testReplicas)
}

func TestCreateUnknownApplicationService(t *testing.T) {
	r, cl := createApplicationReconciler(t, makeTestApplication())
	req := makeRequest()

	res, err := r.Reconcile(req)

	fatalIfError(t, "failed to reconcile", err)
	if res.Requeue != false {
		t.Fatalf("res.Requeue got %v, wanted %v", res.Requeue, false)
	}

	dp := &corev1.Service{}
	err = cl.Get(context.TODO(), ns(testAppName, testNamespace), dp)
	if err != nil {
		t.Fatalf("failed to get created service: %s", err)
	}
}

func TestUpdateExistingConfiguration(t *testing.T) {
	app := makeTestApplication()
	newEnvironment := map[string]string{"new": "value"}
	app.Spec.Environment = newEnvironment
	r, cl := createApplicationReconciler(t, app, configMapFromApplication(makeTestApplication()))
	req := makeRequest()

	_, err := r.Reconcile(req)

	fatalIfError(t, "failed to reconcile", err)
	assertConfigMapHasData(t, testAppName, testNamespace, cl, newEnvironment)
}

func TestUpdateExistingDeployment(t *testing.T) {
	app := makeTestApplication()
	app.Spec.Processes[0].Replicas = 2
	r, cl := createApplicationReconciler(t, app, deploymentFromApplication(makeTestApplication()))
	req := makeRequest()

	_, err := r.Reconcile(req)

	fatalIfError(t, "failed to reconcile", err)
	assertDeploymentConfiguration(t, testAppName, testNamespace, cl, 2)
}

func createApplicationReconciler(t *testing.T, obj ...runtime.Object) (ReconcileApplication, client.Client) {
	scheme := createFakeScheme(t, obj...)
	cl := fake.NewFakeClientWithScheme(scheme, obj...)
	return ReconcileApplication{
		client: cl,
		scheme: scheme,
	}, cl
}

// TODO: Is there a nicer way of registering all types for core and apps?
func createFakeScheme(t *testing.T, objs ...runtime.Object) *runtime.Scheme {
	registerObjs := objs
	registerObjs = append(registerObjs, &corev1.ConfigMap{}, &appsv1.Deployment{}, &corev1.Service{})
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
			Name:      testAppName,
			Namespace: testNamespace,
		},
	}
}

func ns(name, namespace string) types.NamespacedName {
	return types.NamespacedName{Name: name, Namespace: namespace}
}

func assertConfigMapHasData(t *testing.T, name, namespace string, cl client.Client, data map[string]string) {
	t.Helper()
	cm := &corev1.ConfigMap{}
	err := cl.Get(context.TODO(), ns(name+"-config", namespace), cm)
	if err != nil {
		t.Fatalf("failed to get created config-map: %s", err)
	}
	if !reflect.DeepEqual(cm.Data, data) {
		t.Fatalf("got %#v, wanted %#v", cm.Data, data)
	}
}

func assertDeploymentConfiguration(t *testing.T, name, namespace string, cl client.Client, r int32) {
	t.Helper()
	d := &appsv1.Deployment{}
	err := cl.Get(context.TODO(), ns(name, namespace), d)
	if err != nil {
		t.Fatalf("failed to get created deployment: %s", err)
	}
	if *d.Spec.Replicas != r {
		t.Fatalf("got %d, wanted %d", *d.Spec.Replicas, r)
	}
}
