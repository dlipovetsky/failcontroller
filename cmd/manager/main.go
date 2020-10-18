package main

import (
	"os"
	"time"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	controllerruntime "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"daniel.lipovetsky.me/failcontroller"
)

const (
	reconcileTimeout = 10 * time.Second
)

var (
	mgr       manager.Manager
	globalLog = logf.Log.WithName("failcontroller-manager")
)

func main() {
	// controllerruntime.SetLogger(zap.New(zap.UseDevMode(true)))
	controllerruntime.SetLogger(zap.New(zap.UseDevMode(true)))

	mlog := globalLog.WithName("main")
	mlog.Info("loading kubeconfig")
	config := controllerruntime.GetConfigOrDie()

	var err error
	mlog.Info("creating manager")
	mgr, err = controllerruntime.NewManager(config, manager.Options{})
	if err != nil {
		mlog.Error(err, "failed to create manager")
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
			return meta.GetLabels()[failcontroller.KubernetesAppLabel] == failcontroller.Name
		})).
		// Delegate reconciling to reconcileFunc, defined below
		Build(&failcontroller.FailReconciler{
			Manager: mgr,
			Log:     globalLog,
			Timeout: reconcileTimeout,
		}); err != nil {
		mlog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	// Start the Controller through the manager.
	mlog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		mlog.Error(err, "failed to start manager")
		os.Exit(1)
	}
}
