package failcontroller

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type FailReconciler struct {
	Client  client.Client
	Log     logr.Logger
	Timeout time.Duration
}

// The reconciler itself, i.e., code that runs when an Object is
// created/updated/deleted, or processed after a requeue.
func (fr *FailReconciler) Reconcile(r reconcile.Request) (reconcile.Result, error) {
	// Setup logger
	rlog := fr.Log.WithValues("name", r.Name, "namespace", r.Namespace)

	// Setup context
	ctx, cancel := context.WithTimeout(context.Background(), fr.Timeout)
	defer cancel()

	// Get the object being reconciled
	rlog.Info("getting ingress")
	ing := &netv1beta1.Ingress{}
	if err := fr.Client.Get(ctx, r.NamespacedName, ing); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "unable to get ingress")
	}

	if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName != "" {
		// Example of non-idempotent behavior:
		// Every reconcile creates a different ConfigMap.
		if err := fr.Client.Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%d", *ing.Spec.IngressClassName, rand.Intn(100)),
				Namespace: ing.Namespace,
				Labels:    map[string]string{KubernetesAppLabel: Name},
			},
		}); err != nil {
			return reconcile.Result{}, errors.Wrap(err, "unable to create configmap")
		}

		// Example of non-reentrant behavior:
		// The first reconcile creates the ConfigMap. Every subsequent reconcile
		// fails to create the ConfigMap, because it already exists.
		if err := fr.Client.Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      *ing.Spec.IngressClassName,
				Namespace: ing.Namespace,
				Labels:    map[string]string{KubernetesAppLabel: Name},
			},
		}); err != nil {
			return reconcile.Result{}, errors.Wrap(err, "unable to create configmap")
		}
	}
	return reconcile.Result{}, nil
}
