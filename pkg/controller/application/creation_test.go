package application

import (
	"reflect"
	"testing"

	appv1alpha1 "github.com/bigkevmcd/applications/pkg/apis/app/v1alpha1"
)

func TestNewConfigForCR(t *testing.T) {
	config := map[string]string{"testing.value": "42"}

	a := &appv1alpha1.Application{
		Spec: appv1alpha1.ApplicationSpec{
			Replicas: 5,
			Config:   config,
		},
	}

	cm := newConfigMapForCR(a)

	if !reflect.DeepEqual(cm.Data, config) {
		t.Fatalf("newConfigMapForCR() got data %#v, wanted %#v", cm.Data, config)
	}
}
