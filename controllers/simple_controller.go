/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	examplev1 "daniel.lipovetsky.me/failcontroller/api/v1"
)

// SimpleReconciler reconciles a Simple object
type SimpleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.lipovetsky.me,resources=simples,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.lipovetsky.me,resources=simples/status,verbs=get;update;patch

func (r *SimpleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("simple", req.NamespacedName)

	simple := &examplev1.Simple{}
	if err := r.Get(ctx, req.NamespacedName, simple); err != nil {
		if kerrors.IsNotFound(err) {
			log.Info("was deleted")
			return reconcile.Result{}, nil
		}
		log.Error(err, "failed to get")
		return reconcile.Result{}, err
	}

	// Example of non-idempotent behavior:
	// Every reconcile creates a different ConfigMap.
	if err := r.Create(ctx, &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-%d", simple.Name, rand.Intn(100)),
			Namespace:       simple.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(simple, examplev1.GroupVersion.WithKind("Simple"))},
		},
	}); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create configmap")
	}

	// Example of non-reentrant behavior:
	// The first reconcile creates the ConfigMap. Every subsequent reconcile
	// fails to create the ConfigMap, because it already exists.
	if err := r.Create(ctx, &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            simple.Name,
			Namespace:       simple.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(simple, examplev1.GroupVersion.WithKind("Simple"))},
		},
	}); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create configmap")
	}

	return ctrl.Result{}, nil
}

func (r *SimpleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.Simple{}).
		Complete(r)
}
