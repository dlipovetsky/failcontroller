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
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/go-cmp v0.4.1 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/crypto v0.0.0-20200930160638-afb6bcd081ae // indirect
	golang.org/x/net v0.0.0-20200930145003-4acb6c075d10 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c // indirect
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29 // indirect
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800
	sigs.k8s.io/controller-runtime v0.6.3
)
