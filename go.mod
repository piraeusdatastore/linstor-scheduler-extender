module github.com/piraeusdatastore/linstor-scheduler-extender

go 1.17

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
	k8s.io/api => k8s.io/api v0.21.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.4
	k8s.io/apiserver => k8s.io/apiserver v0.21.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.4
	k8s.io/client-go => k8s.io/client-go v0.21.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.4
	k8s.io/code-generator => k8s.io/code-generator v0.21.4
	k8s.io/component-base => k8s.io/component-base v0.21.4
	k8s.io/component-helpers => k8s.io/component-helpers v0.21.4
	k8s.io/controller-manager => k8s.io/controller-manager v0.21.4
	k8s.io/cri-api => k8s.io/cri-api v0.21.4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.4
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.4.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.4
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.4
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.4
	k8s.io/kubectl => k8s.io/kubectl v0.21.4
	k8s.io/kubelet => k8s.io/kubelet v0.21.4
	k8s.io/kubernetes => k8s.io/kubernetes v1.20.11
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.20.4
	k8s.io/metrics => k8s.io/metrics v0.21.4
	k8s.io/mount-utils => k8s.io/mount-utils v0.21.4
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.4
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.9.0
	sigs.k8s.io/sig-storage-lib-external-provisioner/v6 => sigs.k8s.io/sig-storage-lib-external-provisioner/v6 v6.3.0
)

require (
	github.com/LINBIT/golinstor v0.27.0
	github.com/kubernetes-incubator/external-storage v0.20.4-openstorage-rc7
	github.com/libopenstorage/stork v1.4.1-0.20220512171133-b99428ee1ddf
	github.com/portworx/sched-ops v1.20.4-rc1.0.20220401024625-dbc61a336f65
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli v1.22.5
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/component-helpers v0.21.4
	k8s.io/kubernetes v1.21.4
	sigs.k8s.io/controller-runtime v0.9.0
)

require (
	cloud.google.com/go v0.65.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.13 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.8 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.0 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/GoogleCloudPlatform/k8s-cloud-provider v0.0.0-20200415212048-7901bc822317 // indirect
	github.com/aws/aws-sdk-go v1.35.37 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/container-storage-interface/spec v1.5.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/donovanhide/eventsource v0.0.0-20210830082556-c59027999da0 // indirect
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/logr v1.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kubernetes-sigs/aws-ebs-csi-driver v0.9.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/miekg/dns v1.1.35 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.opencensus.io v0.22.4 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5 // indirect
	golang.org/x/sys v0.0.0-20220222200937-f2425489ef4c // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/api v0.30.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210921142501-181ce0d877f6 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apiextensions-apiserver v0.21.4 // indirect
	k8s.io/component-base v0.21.4 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.20.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7 // indirect
	k8s.io/kube-scheduler v0.0.0 // indirect
	k8s.io/utils v0.0.0-20210527160623-6fdb442a123b // indirect
	sigs.k8s.io/gcp-compute-persistent-disk-csi-driver v0.7.0 // indirect
	sigs.k8s.io/sig-storage-lib-external-provisioner/v6 v6.3.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
