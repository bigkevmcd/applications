package application

import (
	"reflect"
	"testing"

	appv1alpha1 "github.com/bigkevmcd/applications/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var testLabels = map[string]string{"app": "my-test-app", "component": "testing"}

func TestConfigMapFromApplication(t *testing.T) {
	config := map[string]string{"testing.value": "42"}
	a := &appv1alpha1.Application{
		Spec: appv1alpha1.ApplicationSpec{
			Labels: testLabels,
			Config: config,
		},
	}

	cm := configMapFromApplication(a)

	if !reflect.DeepEqual(cm.Data, config) {
		t.Fatalf("ConfigMap got data %#v, wanted %#v", cm.Data, config)
	}
	if !reflect.DeepEqual(cm.Labels, testLabels) {
		t.Fatalf("ConfigMap got labels %#v, wanted %#v", cm.Labels, testLabels)
	}
}

func TestDeploymentFromApplication(t *testing.T) {
	containers := []corev1.Container{
		{
			Name:  "nginx",
			Image: "nginx:1.17.4",
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: 80,
				},
			},
		},
	}

	a := &appv1alpha1.Application{
		Spec: appv1alpha1.ApplicationSpec{
			Labels:     testLabels,
			Replicas:   5,
			Containers: containers,
		},
	}

	dp := deploymentFromApplication(a)

	if *dp.Spec.Replicas != 5 {
		t.Fatalf("Deployment got %d Replicas, wanted 5", *dp.Spec.Replicas)
	}
	if !reflect.DeepEqual(dp.Spec.Selector.MatchLabels, testLabels) {
		t.Fatalf("Deployment got %#v MatchLabels, wanted %#v", dp.Spec.Selector.MatchLabels, testLabels)
	}
	if !reflect.DeepEqual(dp.Labels, testLabels) {
		t.Fatalf("Deployment got labels %#v, wanted %#v", dp.Labels, testLabels)
	}
	if !reflect.DeepEqual(dp.Spec.Template.Spec.Containers, containers) {
		t.Fatalf("Deployment got containers %#v, wanted %#v", dp.Spec.Template.Spec.Containers, containers)
	}
	if !reflect.DeepEqual(dp.Spec.Template.ObjectMeta.Labels, testLabels) {
		t.Fatalf("Deployment got deployment labels %#v, wanted %#v", dp.Spec.Template.ObjectMeta.Labels, testLabels)
	}

}

func TestServiceFromApplication(t *testing.T) {
	a := &appv1alpha1.Application{
		Spec: appv1alpha1.ApplicationSpec{
			Labels: testLabels,
		},
	}

	svc := serviceFromApplication(a)

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
