package experimental

import (
	"context"
	"reflect"
	"strconv"

	experimentsv1alpha1 "experimental-operator/pkg/apis/experiments/v1alpha1"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_experimental")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Experimental Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileExperimental{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("experimental-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Experimental
	err = c.Watch(&source.Kind{Type: &experimentsv1alpha1.Experimental{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Experimental
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &experimentsv1alpha1.Experimental{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileExperimental implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileExperimental{}

// ReconcileExperimental reconciles a Experimental object
type ReconcileExperimental struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Experimental object and makes changes based on the state read
// and what is in the Experimental.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileExperimental) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Experimental")

	// Fetch the Experimental instance
	instance := &experimentsv1alpha1.Experimental{}
	if err := r.client.Get(context.TODO(), request.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(instance, 0)

	// Set Experimental instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)

	// Fetch the PodSet instance
	experimentalResource := &experimentsv1alpha1.Experimental{}

	if err = r.client.Get(context.TODO(), request.NamespacedName, experimentalResource); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	existingPods, existingPodNames, err := r.existingPodInfo(request, experimentalResource, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	return r.reconcileReplicas(request, existingPods, existingPodNames, experimentalResource, reqLogger)
}

func (r *ReconcileExperimental) existingPodInfo(request reconcile.Request, experimentalResource *experimentsv1alpha1.Experimental, logger logr.Logger) (*corev1.PodList, []string, error) {
	// List all pods owned by this Experimental instance
	lbls := labels.Set{
		"app": experimentalResource.Name,
	}

	existingPods := &corev1.PodList{}

	err := r.client.List(context.TODO(),
		existingPods,
		&client.ListOptions{
			Namespace:     request.Namespace,
			LabelSelector: labels.SelectorFromSet(lbls),
		})
	if err != nil {
		logger.Error(err, "failed to list existing pods in the experimental resource")
		return nil, nil, err
	}

	existingPodNames := []string{}

	// Count the pods that are pending or running as available
	for _, pod := range existingPods.Items {
		if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
			continue
		}

		if pod.Status.Phase == corev1.PodPending || pod.Status.Phase == corev1.PodRunning {
			existingPodNames = append(existingPodNames, pod.GetObjectMeta().GetName())
		}
	}

	logger.Info("Checking podset", "expected replicas", experimentalResource.Spec.NumberOfReplicas, "Pod.Names", existingPodNames)

	return existingPods, existingPodNames, nil
}

func (r *ReconcileExperimental) reconcileReplicas(request reconcile.Request, existingPods *corev1.PodList, existingPodNames []string, experimentalResource *experimentsv1alpha1.Experimental, logger logr.Logger) (reconcile.Result, error) {
	// Update the status if necessary
	status := experimentsv1alpha1.ExperimentalStatus{
		NumberOfReplicas:      int32(len(existingPodNames)),
		PersistentVolumesUsed: []string{"This would be one in a list of PVs", "Found by matching against the app label", "And extracting the appropriate metadata"},
	}

	if !reflect.DeepEqual(experimentalResource.Status, status) {
		experimentalResource.Status = status
		if err := r.client.Status().Update(context.TODO(), experimentalResource); err != nil {
			logger.Error(err, "failed to update the experimental resource")
			return reconcile.Result{}, err
		}
	}

	if err := r.scaleDownPods(experimentalResource, existingPods, existingPodNames, logger); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.scaleUpPods(experimentalResource, existingPods, existingPodNames, logger); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{Requeue: true}, nil
}

func (r *ReconcileExperimental) scaleDownPods(experimentalResource *experimentsv1alpha1.Experimental, existingPods *corev1.PodList, existingPodNames []string, logger logr.Logger) error {
	if int32(len(existingPodNames)) > experimentalResource.Spec.NumberOfReplicas {
		// delete a pod. Just one at a time (this reconciler will be called again afterwards)
		logger.Info("Deleting a pod in the experimental resource", "expected replicas", experimentalResource.Spec.NumberOfReplicas, "Pod.Names", existingPodNames)

		pod := existingPods.Items[len(existingPodNames)-1]

		if err := r.client.Delete(context.TODO(), &pod); err != nil {
			logger.Error(err, "failed to delete a pod")
			return err
		}
	}

	return nil
}

func (r *ReconcileExperimental) scaleUpPods(experimentalResource *experimentsv1alpha1.Experimental, existingPods *corev1.PodList, existingPodNames []string, logger logr.Logger) error {
	if int32(len(existingPodNames)) < experimentalResource.Spec.NumberOfReplicas {
		// create a new pod. Just one at a time (this reconciler will be called again afterwards)
		logger.Info("Adding a pod in the experimentalResource", "expected replicas", experimentalResource.Spec.NumberOfReplicas, "Pod.Names", existingPodNames)

		pod := newPodForCR(experimentalResource, len(existingPodNames))

		if err := controllerutil.SetControllerReference(experimentalResource, pod, r.scheme); err != nil {
			logger.Error(err, "unable to set owner reference on new pod")
			return err
		}

		err := r.client.Create(context.TODO(), pod)
		if err != nil {
			logger.Error(err, "failed to create a pod")
			return err
		}
	}

	return nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *experimentsv1alpha1.Experimental, replicaCount int) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod" + strconv.Itoa(replicaCount),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
