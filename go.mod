module github.com/piraeusdatastore/linstor-scheduler-extender

go 1.19

replace (
	github.com/LINBIT/golinstor => github.com/LINBIT/golinstor v0.41.2
	github.com/banzaicloud/k8s-objectmatcher => github.com/banzaicloud/k8s-objectmatcher v1.5.1
	github.com/census-instrumentation/opencensus-proto => github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/docker/distribution => github.com/docker/distribution v2.7.0+incompatible
	github.com/docker/docker => github.com/moby/moby v20.10.4+incompatible
	github.com/go-logr/logr => github.com/go-logr/logr v0.3.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1
	github.com/heptio/velero => github.com/heptio/velero v1.0.0
	github.com/kubernetes-csi/external-snapshotter/client/v4 => github.com/kubernetes-csi/external-snapshotter/client/v4 v4.0.0
	github.com/kubernetes-incubator/external-storage => github.com/libopenstorage/external-storage v0.20.4-openstorage-rc7
	github.com/kubernetes-incubator/external-storage v0.20.4-openstorage-rc7 => github.com/libopenstorage/external-storage v0.20.4-openstorage-rc7
	github.com/libopenstorage/autopilot-api => github.com/libopenstorage/autopilot-api v0.6.1-0.20210301232050-ca2633c6e114
	github.com/portworx/sched-ops => github.com/portworx/sched-ops v1.20.4-rc1.0.20220401024625-dbc61a336f65
	github.com/portworx/torpedo => github.com/portworx/torpedo v0.0.0-20220413175518-a912f1a6faa2

	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.8.2
	gopkg.in/fsnotify.v1 v1.4.7 => github.com/fsnotify/fsnotify v1.4.7
	helm.sh/helm/v3 => helm.sh/helm/v3 v3.6.0
	k8s.io/api => k8s.io/api v0.25.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.25.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.25.6
	k8s.io/apiserver => k8s.io/apiserver v0.25.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.25.6
	k8s.io/client-go => k8s.io/client-go v0.25.6
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.25.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.25.6
	k8s.io/code-generator => k8s.io/code-generator v0.25.6
	k8s.io/component-base => k8s.io/component-base v0.25.6
	k8s.io/component-helpers => k8s.io/component-helpers v0.25.6
	k8s.io/controller-manager => k8s.io/controller-manager v0.25.6
	k8s.io/cri-api => k8s.io/cri-api v0.25.6
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.25.6
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.4.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.25.6
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.25.6
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.25.6
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.25.6
	k8s.io/kubectl => k8s.io/kubectl v0.25.6
	k8s.io/kubelet => k8s.io/kubelet v0.25.6
	k8s.io/kubernetes => k8s.io/kubernetes v1.25.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.25.6
	k8s.io/metrics => k8s.io/metrics v0.25.6
	k8s.io/mount-utils => k8s.io/mount-utils v0.25.6
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.25.6
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.25.6
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.9.0
	sigs.k8s.io/sig-storage-lib-external-provisioner/v6 => sigs.k8s.io/sig-storage-lib-external-provisioner/v6 v6.3.0
)

require (
	github.com/LINBIT/golinstor v0.46.1
	github.com/kubernetes-incubator/external-storage v0.25.1-openstorage-rc1
	github.com/libopenstorage/stork v2.6.4-rc2+incompatible
	github.com/portworx/sched-ops v1.20.4-rc1.0.20220401024625-dbc61a336f65
	github.com/sirupsen/logrus v1.9.3
	github.com/slok/kubewebhook/v2 v2.5.0
	github.com/urfave/cli v1.22.14
	k8s.io/api v0.25.6
	k8s.io/apimachinery v0.25.6
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/component-helpers v0.25.6
	sigs.k8s.io/controller-runtime v0.9.0
)

require (
	cloud.google.com/go/compute v1.19.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.27 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.22 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/donovanhide/eventsource v0.0.0-20210830082556-c59027999da0 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/glog v1.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/miekg/dns v1.1.35 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/oauth2 v0.7.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/term v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	gomodules.xyz/jsonpatch/v3 v3.0.1 // indirect
	gomodules.xyz/orderedmap v0.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.25.6 // indirect
	k8s.io/component-base v0.25.6 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1 // indirect
	k8s.io/kubernetes v1.25.1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/sig-storage-lib-external-provisioner v4.1.0+incompatible // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
