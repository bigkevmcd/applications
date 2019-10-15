package application

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appv1alpha1 "github.com/bigkevmcd/applications/pkg/apis/app/v1alpha1"
)

var (
	testLabels      = map[string]string{"app.kubernetes.io/name": testAppName}
	testEnvironment = map[string]string{"TEST_MODE": "true"}
	testImage       = "test-image:latest"
	testProcess     = appv1alpha1.ProcessSpec{
		Name:     "web",
		Replicas: 5,
		Image:    testImage,
		Port:     80,
	}
)

func TestConfigMapFromApplication(t *testing.T) {
	app := makeTestApplication()

	cm := configMapFromApplication(app)

	if !reflect.DeepEqual(cm.Data, testEnvironment) {
		t.Fatalf("ConfigMap got data %#v, wanted %#v", cm.Data, testEnvironment)
	}
	if !reflect.DeepEqual(cm.Labels, testLabels) {
		t.Fatalf("ConfigMap got labels %#v, wanted %#v", cm.Labels, testLabels)
	}
}

func TestMakeEnvFromApp(t *testing.T) {
	testMode := "TEST_MODE"
	app := makeTestApplication()

	env := makeEnvFromApp(app)

	if l := len(env); l != 1 {
		t.Fatalf("makeEnvFromApp() got %d vars, wanted 1", l)
	}
	v := env[0]

	wanted := corev1.EnvVar{
		Name: testMode,
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMapNameForApp(app),
				},
				Key: "TEST_MODE",
			},
		},
	}
	if !reflect.DeepEqual(v, wanted) {
		t.Fatalf("makeEnvFromApp() got %#v, wanted %#v", v, wanted)
	}

}

func TestMakePodSpec(t *testing.T) {
	app := makeTestApplication()

	s := makePodSpec(app, testProcess)

	wanted := corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  app.ObjectMeta.Name + "-" + "web",
				Image: testImage,
				Env:   makeEnvFromApp(app),
			},
		},
	}
	if !reflect.DeepEqual(s, wanted) {
		t.Fatalf("makePodSpec() got %#v, wanted %#v", s, wanted)
	}
}

func TestDeploymentFromApplication(t *testing.T) {
	process := appv1alpha1.ProcessSpec{
		Name:     "web",
		Replicas: 5,
		Image:    testImage,
		Port:     80,
	}
	app := makeTestApplication()
	app.Spec.Processes = []appv1alpha1.ProcessSpec{process}

	dp := deploymentFromApplication(app)

	if *dp.Spec.Replicas != 5 {
		t.Fatalf("Deployment got %d Replicas, wanted 5", *dp.Spec.Replicas)
	}
	if !reflect.DeepEqual(dp.Spec.Selector.MatchLabels, testLabels) {
		t.Fatalf("Deployment got %#v MatchLabels, wanted %#v", dp.Spec.Selector.MatchLabels, testLabels)
	}
	if !reflect.DeepEqual(dp.Labels, testLabels) {
		t.Fatalf("Deployment got labels %#v, wanted %#v", dp.Labels, testLabels)
	}
	if l := len(dp.Spec.Template.Spec.Containers); l != 1 {
		t.Fatalf("Deployment got %d containers, wanted 1", l)
	}

	wantedContainer := corev1.Container{
		Name:  app.ObjectMeta.Name + "-web",
		Image: testImage,
		Env:   makeEnvFromApp(app),
	}
	if !reflect.DeepEqual(dp.Spec.Template.Spec.Containers[0], wantedContainer) {
		t.Fatalf("Deployment got containers %#v, wanted %#v", dp.Spec.Template.Spec.Containers[0], wantedContainer)
	}
	if !reflect.DeepEqual(dp.Spec.Template.ObjectMeta.Labels, testLabels) {
		t.Fatalf("Deployment got deployment labels %#v, wanted %#v", dp.Spec.Template.ObjectMeta.Labels, testLabels)
	}

}

func TestServiceFromApplication(t *testing.T) {
	app := makeTestApplication()

	svc := serviceFromApplication(app)

	wanted := []corev1.ServicePort{
		{
			Protocol: corev1.ProtocolTCP,
			Port:     80,
		},
	}
	if !reflect.DeepEqual(svc.Spec.Ports, wanted) {
		t.Fatalf("Service got ports %#v, wanted %#v", svc.Spec.Ports, wanted)
	}

	if !reflect.DeepEqual(svc.Spec.Selector, testLabels) {
		t.Fatalf("Service got selector %#v, wanted %#v", svc.Spec.Selector, testLabels)
	}
	if svc.Spec.Type != corev1.ServiceTypeNodePort {
		t.Fatalf("Service got type %s, wanted %s", svc.Spec.Type, corev1.ServiceTypeNodePort)
	}
}

func makeTestApplication() *appv1alpha1.Application {
	return &appv1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testAppName,
			Namespace: testNamespace,
		},
		Spec: appv1alpha1.ApplicationSpec{
			Environment: testEnvironment,
			Processes:   []appv1alpha1.ProcessSpec{testProcess},
		},
	}
}
