package application

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-logr/logr"

	appv1alpha1 "github.com/bigkevmcd/applications/pkg/apis/app/v1alpha1"
)

var log = logf.Log.WithName("controller_application")

// Add creates a new Application Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileApplication{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler.
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("application-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appv1alpha1.Application{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO: find out how to test this.
	watchedTypes := []runtime.Object{&corev1.ConfigMap{}, &appsv1.Deployment{}, &corev1.Service{}}
	for _, t := range watchedTypes {
		err = c.Watch(&source.Kind{Type: t}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &appv1alpha1.Application{},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ReconcileApplication reconciles an Application object.
type ReconcileApplication struct {
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Application object and makes
// changes based on the state read and what is in the Application.Spec.
func (r *ReconcileApplication) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Application")

	application := &appv1alpha1.Application{}
	err := r.client.Get(context.TODO(), request.NamespacedName, application)

	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	err = r.createOrUpdateConfigMap(application, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.createOrUpdateDeployment(application, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.createOrUpdateService(application, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, err
}

func (r *ReconcileApplication) createOrUpdateConfigMap(a *appv1alpha1.Application, logger logr.Logger) error {
	configMap := configMapFromApplication(a)
	err := controllerutil.SetControllerReference(a, configMap, r.scheme)
	if err != nil {
		return err
	}

	found := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new ConfigMap", "Created.Namespace", configMap.Namespace, "Created.Name", configMap.Name)
		err = r.client.Create(context.TODO(), configMap)
		if err != nil {
			return err
		}
		a.Status.ConfigMapName = configMap.Name
		return r.client.Status().Update(context.TODO(), a)
	} else if err != nil {
		return err
	}

	logger.Info("Updating existing ConfigMap", "Updated.Namespace", configMap.Namespace, "Updated.Name", configMap.Name)
	found.Data = configMap.Data
	err = r.client.Update(context.TODO(), found)
	if err != nil {
		return err
	}
	a.Status.ConfigMapName = configMap.Name
	return r.client.Status().Update(context.TODO(), a)
}

func (r *ReconcileApplication) createOrUpdateDeployment(a *appv1alpha1.Application, logger logr.Logger) error {
	deployment := deploymentFromApplication(a)
	err := controllerutil.SetControllerReference(a, deployment, r.scheme)
	if err != nil {
		return err
	}

	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Deployment", "Created.Namespace", deployment.Namespace, "Created.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return err
		}
		a.Status.DeploymentName = deployment.Name
		return r.client.Status().Update(context.TODO(), a)
	} else if err != nil {
		return err
	}

	logger.Info("Updating existing Deployment", "Updated.Namespace", deployment.Namespace, "Updated.Name", deployment.Name)
	found.Spec = deployment.Spec
	err = r.client.Update(context.TODO(), found)
	if err != nil {
		return err
	}
	a.Status.ConfigMapName = deployment.Name
	return r.client.Status().Update(context.TODO(), a)
}

func (r *ReconcileApplication) createOrUpdateService(a *appv1alpha1.Application, logger logr.Logger) error {
	service := serviceFromApplication(a)
	err := controllerutil.SetControllerReference(a, service, r.scheme)
	if err != nil {
		return err
	}

	found := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Service", "Created.Namespace", service.Namespace, "Created.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return err
		}
		a.Status.ServiceName = service.Name
		return r.client.Status().Update(context.TODO(), a)
	} else if err != nil {
		return err
	}

	logger.Info("Updating existing Service", "Updated.Namespace", service.Namespace, "Updated.Name", service.Name)
	found.Spec = service.Spec
	err = r.client.Update(context.TODO(), found)
	if err != nil {
		return err
	}
	a.Status.ServiceName = service.Name
	return r.client.Status().Update(context.TODO(), a)
}
