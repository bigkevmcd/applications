package application

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appv1alpha1 "github.com/bigkevmcd/applications/pkg/apis/app/v1alpha1"
)

// configMapFromApplication makes a ConfigMap based on the Application.
func configMapFromApplication(app *appv1alpha1.Application) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: makeObjectMeta(configMapNameForApp(app), app),
		Data:       app.Spec.Environment,
	}
}

// deploymentFromApplication makes a deployment based on the Application.
// TODO: fix this for multiple processes.
func deploymentFromApplication(app *appv1alpha1.Application) *appsv1.Deployment {
	process := app.Spec.Processes[0]
	return &appsv1.Deployment{
		ObjectMeta: makeObjectMeta(app.Name, app),
		Spec: appsv1.DeploymentSpec{
			Replicas: &process.Replicas,
			Selector: makeLabelSelector(app),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: makeObjectMeta(app.Name, app),
				Spec:       makePodSpec(app, process),
			},
		},
	}
}

// serviceFromApplication makes a service based on the Application.
// TODO: What to do about configuring the service type, port and protocol?
func serviceFromApplication(app *appv1alpha1.Application) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: makeObjectMeta(app.Name, app),
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: labelsForApp(app),
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
				},
			},
		},
	}
}

func makeObjectMeta(name string, app *appv1alpha1.Application) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: app.Namespace,
		Labels:    labelsForApp(app),
	}
}

func makeLabelSelector(app *appv1alpha1.Application) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: labelsForApp(app),
	}
}

func labelsForApp(app *appv1alpha1.Application) map[string]string {
	return map[string]string{"app": app.ObjectMeta.Name}
}

func makePodSpec(app *appv1alpha1.Application, p appv1alpha1.ProcessSpec) corev1.PodSpec {
	return corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  app.ObjectMeta.Name + "-" + p.Name,
				Image: p.Image,
				Env:   makeEnvFromApp(app),
			},
		},
	}
}

func makeEnvFromApp(app *appv1alpha1.Application) []corev1.EnvVar {
	vars := []corev1.EnvVar{}
	for k, _ := range app.Spec.Environment {
		envVar := corev1.EnvVar{
			Name: k,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configMapNameForApp(app),
					},
					Key: k,
				},
			},
		}
		vars = append(vars, envVar)
	}
	return vars
}

func configMapNameForApp(app *appv1alpha1.Application) string {
	return app.Name + "-config"
}
