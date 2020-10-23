package failcontroller

import (
	"context"
	"time"

	kpointer "k8s.io/utils/pointer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	// corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	reconcileTimeout = 10 * time.Second
)

var _ = Describe("AWSClusterReconciler", func() {
	BeforeEach(func() {})
	AfterEach(func() {})

	Context("Reconcile an AWSCluster", func() {
		It("should not error and not requeue the request with insufficient set up", func() {
			ctx := context.Background()

			reconciler := &FailReconciler{
				Client: k8sClient,
				// XXX: null logger
				Log:     log.Log,
				Timeout: reconcileTimeout,
			}

			prefixPathType := extv1beta1.PathTypePrefix

			// Create the Ingress object
			ing := &extv1beta1.Ingress{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "networking.k8s.io/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: "default",
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "failcontroller",
					},
				},
				Spec: extv1beta1.IngressSpec{
					IngressClassName: kpointer.StringPtr("foo"),
					Rules: []extv1beta1.IngressRule{
						{
							IngressRuleValue: extv1beta1.IngressRuleValue{
								HTTP: &extv1beta1.HTTPIngressRuleValue{
									Paths: []extv1beta1.HTTPIngressPath{
										{
											Path:     "/testpath",
											PathType: &prefixPathType,
											Backend: extv1beta1.IngressBackend{
												ServiceName: "test",
												ServicePort: intstr.IntOrString{
													IntVal: 80,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, ing)).To(Succeed())

			result, err := reconciler.Reconcile(ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: ing.Namespace,
					Name:      ing.Name,
				},
			})
			Expect(err).To(BeNil())
			Expect(result.RequeueAfter).To(BeZero())
		})
	})
})
