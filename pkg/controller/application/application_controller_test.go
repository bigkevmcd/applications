package application

import (
	"context"
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
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
	testNamespace = "testing"
	testAppName   = "test-application"
	testReplicas  = 5
)

var (
	testConfig = map[string]string{"testing": "value"}
)

func TestCreateUnknownApplicationConfiguration(t *testing.T) {
	r, cl := createApplicationReconciler(t, makeApplication())
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
	if wanted := testAppName + "-config"; wanted != app.Status.ConfigMapName {
		t.Fatalf("got %s, wanted %v", app.Status, wanted)
	}

	assertConfigMapHasData(t, testAppName, testNamespace, cl, testConfig)
}

func TestCreateUnknownApplicationDeployment(t *testing.T) {
	r, cl := createApplicationReconciler(t, makeApplication())
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

func TestUpdateExistingConfiguration(t *testing.T) {
	app := makeApplication()
	newConfig := map[string]string{"new": "value"}
	app.Spec.Config = newConfig
	r, cl := createApplicationReconciler(t, app, configMapFromApplication(makeApplication()))
	req := makeRequest()

	_, err := r.Reconcile(req)

	fatalIfError(t, "failed to reconcile", err)
	assertConfigMapHasData(t, testAppName, testNamespace, cl, newConfig)
}

func TestUpdateExistingDeployment(t *testing.T) {
	app := makeApplication()
	app.Spec.Replicas = 2
	r, cl := createApplicationReconciler(t, app, deploymentFromApplication(makeApplication()))
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

func createFakeScheme(t *testing.T, objs ...runtime.Object) *runtime.Scheme {
	registerObjs := objs
	registerObjs = append(registerObjs, &corev1.ConfigMap{}, &appsv1.Deployment{})
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

func makeApplication() *api.Application {
	return &api.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testAppName,
			Namespace: testNamespace,
		},
		Spec: api.ApplicationSpec{
			Config:   testConfig,
			Replicas: testReplicas,
		},
	}
}

func ns(name, namespace string) types.NamespacedName {
	return types.NamespacedName{Name: name, Namespace: namespace}
}

func assertConfigMapHasData(t *testing.T, name, namespace string, cl client.Client, data map[string]string) {
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
	d := &appsv1.Deployment{}
	err := cl.Get(context.TODO(), ns(name, namespace), d)
	if err != nil {
		t.Fatalf("failed to get created deployment: %s", err)
	}
	if *d.Spec.Replicas != r {
		t.Fatalf("got %d, wanted %d", *d.Spec.Replicas, r)
	}
}
