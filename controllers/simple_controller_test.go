package controllers

import (
	"context"
	"time"

	examplev1 "daniel.lipovetsky.me/failcontroller/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var _ = Describe("The controller", func() {
	It("should be reentrant", func() {
		ctx := context.TODO()

		stopCh := make(chan struct{})
		mgr, err := ctrl.NewManager(cfg, ctrl.Options{})
		Expect(err).NotTo(HaveOccurred(), "failed to create manager")

		controller := &SimpleReconciler{
			Client: mgr.GetClient(),
			Log:    logf.Log,
		}
		err = controller.SetupWithManager(mgr)
		Expect(err).NotTo(HaveOccurred(), "failed to setup controller")

		go func() {
			err := mgr.Start(stopCh)
			Expect(err).NotTo(HaveOccurred(), "failed to start manager")
		}()

		wantObjectKey := client.ObjectKey{
			Name:      "testresource",
			Namespace: metav1.NamespaceDefault,
		}
		want := &examplev1.Simple{
			ObjectMeta: metav1.ObjectMeta{
				Name:      wantObjectKey.Name,
				Namespace: wantObjectKey.Namespace,
			},
			Spec: examplev1.SimpleSpec{
				Foo: "foo",
			},
		}
		err = k8sClient.Create(ctx, want)
		Expect(err).NotTo(HaveOccurred())

		got := &examplev1.Simple{}
		Eventually(
			getResourceFunc(ctx, wantObjectKey, got),
			5*time.Second,
			500*time.Millisecond,
		).Should(BeNil(), "simple resource should exist")

		close(stopCh)
	})
})

func getResourceFunc(ctx context.Context, key client.ObjectKey, obj runtime.Object) func() error {
	return func() error {
		return k8sClient.Get(ctx, key, obj)
	}
}
