module daniel.lipovetsky.me/failcontroller

go 1.15

replace (
	k8s.io/api => k8s.io/api v0.18.8
	k8s.io/apiserver => k8s.io/apiserver v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
	k8s.io/kubectl => k8s.io/kubectl v0.18.8
)

require (
	github.com/go-logr/logr v0.1.0
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	sigs.k8s.io/controller-runtime v0.6.3
)
