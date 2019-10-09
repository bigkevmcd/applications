package application

import (
	"reflect"
	"testing"

	appv1alpha1 "github.com/bigkevmcd/applications/pkg/apis/app/v1alpha1"
	"k8s.io/api/core/v1"
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
		t.Fatalf("configMapFromApplication() got data %#v, wanted %#v", cm.Data, config)
	}
	if !reflect.DeepEqual(cm.Labels, testLabels) {
		t.Fatalf("configMapFromApplication() got labels %#v, wanted %#v", cm.Labels, testLabels)
	}
}

func TestDeploymentFromApplication(t *testing.T) {
	a := &appv1alpha1.Application{
		Spec: appv1alpha1.ApplicationSpec{
			Labels:     testLabels,
			Replicas:   5,
			Containers: []*v1.Container{},
		},
	}

	dp := deploymentFromApplication(a)

	if *dp.Spec.Replicas != 5 {
		t.Fatalf("deploymentFromApplication() got %d Replicas, wanted 5", *dp.Spec.Replicas)
	}
	if !reflect.DeepEqual(dp.Labels, testLabels) {
		t.Fatalf("configMapFromApplication() got labels %#v, wanted %#v", dp.Labels, testLabels)
	}
}
