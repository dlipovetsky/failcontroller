package failcontroller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	kpointer "k8s.io/utils/pointer"
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
		It("should ", func() {
			ctx := context.Background()

			reconciler := &FailReconciler{
				Client: k8sClient,
				// XXX: null logger
				Log:     log.Log,
				Timeout: reconcileTimeout,
			}

			prefixPathType := netv1beta1.PathTypePrefix

			// Create the Ingress object
			ing := &netv1beta1.Ingress{
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
				Spec: netv1beta1.IngressSpec{
					IngressClassName: kpointer.StringPtr("fooclass"),
					Rules: []netv1beta1.IngressRule{
						{
							IngressRuleValue: netv1beta1.IngressRuleValue{
								HTTP: &netv1beta1.HTTPIngressRuleValue{
									Paths: []netv1beta1.HTTPIngressPath{
										{
											Path:     "/testpath",
											PathType: &prefixPathType,
											Backend: netv1beta1.IngressBackend{
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
					Name:      ing.Name,
					Namespace: ing.Namespace,
				},
			})
			Expect(err).To(BeNil())
			Expect(result.RequeueAfter).To(BeZero())

			want := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      *ing.Spec.IngressClassName,
					Namespace: ing.Namespace,
					Labels:    map[string]string{KubernetesAppLabel: Name},
				},
			}

			got := &corev1.ConfigMap{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: want.Name, Namespace: want.Namespace}, got)).To(Succeed())
		})
	})
})
