package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	controllerruntime "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	annotationKey    = "lipovetsky.daniel.me/name"
	annotationValue  = "controller-example"
	reconcileTimeout = 10 * time.Second
)

var (
	mgr       manager.Manager
	globalLog = logf.Log.WithName("controller-example")
)

func main() {
	controllerruntime.SetLogger(zap.New(zap.UseDevMode(true)))

	mlog := globalLog.WithName("main")
	mlog.Info("loading kubeconfig")
	config := controllerruntime.GetConfigOrDie()

	var err error
	mlog.Info("creating manager")
	mgr, err = controllerruntime.NewManager(config, manager.Options{})
	if err != nil {
		mlog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Create the controller, and register it with the manager
	mlog.Info("creating controller")
	if _, err := controllerruntime.NewControllerManagedBy(mgr).
		// Reconcile v1beta1.Ingress.
		For(&v1beta1.Ingress{}).
		// But include only those with the right annotation
		// (Server-side filtering will be available in the future, see https://github.com/kubernetes-sigs/controller-runtime/issues/244)
		WithEventFilter(predicate.NewPredicateFuncs(func(meta metav1.Object, object runtime.Object) bool {
			return meta.GetAnnotations()[annotationKey] == annotationValue
		})).
		// Delegate reconciling to reconcileFunc, defined below
		Build(reconcile.Func(reconcileFunc)); err != nil {
		mlog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	// Start the Controller through the manager.
	mlog.Info("continuing to run manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		mlog.Error(err, "unable to continue running manager")
		os.Exit(1)
	}
}

// The reconciler itself, i.e., code that runs when an Object is
// created/updated/deleted, or processed after a requeue.
func reconcileFunc(r reconcile.Request) (reconcile.Result, error) {
	// Setup logger
	rlog := globalLog.WithValues("kind", v1beta1.Ingress{}.Kind, "name", r.Name, "namespace", r.Namespace)

	// Setup context
	ctx, cancel := context.WithTimeout(context.Background(), reconcileTimeout)
	defer cancel()

	// Setup API client
	c := mgr.GetClient()

	// Get the object being reconciled
	rlog.Info("getting object")
	ing := &v1beta1.Ingress{}
	if err := c.Get(ctx, r.NamespacedName, ing); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "unable to get object")
	}

	rlog.Info("updating status")
	ing.Status.LoadBalancer.Ingress = []v1.LoadBalancerIngress{
		{
			Hostname: fmt.Sprintf("host%d.example.com", rand.Int()),
		},
	}
	if err := c.Status().Update(ctx, ing); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "unable to update status")
	}

	return reconcile.Result{}, nil
}
