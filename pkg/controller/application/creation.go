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
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-config",
			Namespace: cr.Namespace,
			Labels:    cr.Spec.Labels,
		},
		Data: cr.Spec.Config,
	}
}

// deploymentFromApplication makes a deployment based on the Application.
func deploymentFromApplication(cr *appv1alpha1.Application) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    cr.Spec.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cr.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: cr.Spec.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name,
					Namespace: cr.Namespace,
					Labels:    cr.Spec.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: cr.Spec.Containers,
				},
			},
		},
	}
}

// serviceFromApplication makes a service based on the Application.
func serviceFromApplication(cr *appv1alpha1.Application) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    cr.Spec.Labels,
		},
	}
}
