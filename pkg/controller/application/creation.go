package application

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appv1alpha1 "github.com/bigkevmcd/applications/pkg/apis/app/v1alpha1"
)

// TODO: What should this do if we get no labels, autogenerate them based on
// the Name?

// configMapFromApplication makes a ConfigMap based on the Application.
func configMapFromApplication(cr *appv1alpha1.Application) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: makeObjectMeta(cr.Name+"-config", cr),
		Data:       cr.Spec.Config,
	}
}

// deploymentFromApplication makes a deployment based on the Application.
func deploymentFromApplication(cr *appv1alpha1.Application) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: makeObjectMeta(cr.Name, cr),
		Spec: appsv1.DeploymentSpec{
			Replicas: &cr.Spec.Replicas,
			Selector: makeLabelSelector(cr),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: makeObjectMeta(cr.Name, cr),
				Spec: corev1.PodSpec{
					Containers: cr.Spec.Containers,
				},
			},
		},
	}
}

// serviceFromApplication makes a service based on the Application.
// TODO: What to do about configuring the service type, port and protocol?
func serviceFromApplication(cr *appv1alpha1.Application) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: makeObjectMeta(cr.Name, cr),
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: cr.Spec.Labels,
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
				},
			},
		},
	}
}

func makeObjectMeta(name string, cr *appv1alpha1.Application) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: cr.Namespace,
		Labels:    cr.Spec.Labels,
	}
}

func makeLabelSelector(cr *appv1alpha1.Application) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: cr.Spec.Labels,
	}
}
