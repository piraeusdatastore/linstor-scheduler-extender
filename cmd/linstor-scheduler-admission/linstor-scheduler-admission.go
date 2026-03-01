package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const annBetaStorageProvisioner = "volume.beta.kubernetes.io/storage-provisioner"
const annStorageProvisioner = "volume.kubernetes.io/storage-provisioner"

type config struct {
	certFile      string
	keyFile       string
	driverName    string
	schedulerName string
}

func initFlags() *config {
	cfg := &config{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.certFile, "tls-cert-file", "", "TLS certificate file")
	fl.StringVar(&cfg.keyFile, "tls-key-file", "", "TLS key file")
	fl.StringVar(&cfg.driverName, "driver", "linstor.csi.linbit.com", "Driver name")
	fl.StringVar(&cfg.schedulerName, "scheduler", "linstor", "Scheduler name")

	fl.Parse(os.Args[1:])
	return cfg
}

type jsonPatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// shouldMutateScheduler checks if a pod uses LINSTOR volumes and should have its scheduler overridden.
func shouldMutateScheduler(ctx context.Context, pod *corev1.Pod, namespace string, cli kubernetes.Interface, cfg *config) bool {
	var pvcNames []string

	for _, volume := range pod.Spec.Volumes {
		// Volume has inline CSI driver assigned
		if volume.CSI != nil && volume.CSI.Driver == cfg.driverName {
			return true
		}
		// Volume is not PVC, it does not interest us
		if volume.PersistentVolumeClaim == nil {
			continue
		}
		pvcNames = append(pvcNames, volume.PersistentVolumeClaim.ClaimName)
	}

	// Check PVCs
	for _, pvcName := range pvcNames {
		var discoveredProvisioner string
		pvc, err := cli.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			logrus.Warnf("Failed to get PVC %s/%s: %v", namespace, pvcName, err)
			continue
		}
		// Try to gather provisioner name from annotations
		if pvc != nil {
			if provisioner, ok := pvc.Annotations[annStorageProvisioner]; ok {
				discoveredProvisioner = provisioner
			}
			if provisioner, ok := pvc.Annotations[annBetaStorageProvisioner]; ok {
				discoveredProvisioner = provisioner
			}
		}
		// Try to gather provisioner name from associated StorageClass
		if discoveredProvisioner == "" && pvc != nil && pvc.Spec.StorageClassName != nil && *pvc.Spec.StorageClassName != "" {
			sc, err := cli.StorageV1().StorageClasses().Get(ctx, *pvc.Spec.StorageClassName, metav1.GetOptions{})
			if err != nil && !errors.IsNotFound(err) {
				logrus.Warnf("Failed to get StorageClass %s: %v", *pvc.Spec.StorageClassName, err)
				continue
			}
			if sc != nil && sc.Provisioner == cfg.driverName {
				discoveredProvisioner = sc.Provisioner
			}
		}
		// Try to gather provisioner name from associated PV
		if discoveredProvisioner == "" && pvc != nil && pvc.Spec.VolumeName != "" {
			pv, err := cli.CoreV1().PersistentVolumes().Get(ctx, pvc.Spec.VolumeName, metav1.GetOptions{})
			if err != nil && !errors.IsNotFound(err) {
				logrus.Warnf("Failed to get PV %s: %v", pvc.Spec.VolumeName, err)
				continue
			}
			if pv != nil && pv.Spec.CSI != nil {
				discoveredProvisioner = pv.Spec.CSI.Driver
			}
		}
		// Overwrite the scheduler name
		if discoveredProvisioner == cfg.driverName {
			return true
		}
	}

	return false
}

func handleMutate(w http.ResponseWriter, r *http.Request, cli kubernetes.Interface, cfg *config) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Failed to read request body: %v", err)
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var review admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &review); err != nil {
		logrus.Errorf("Failed to decode admission review: %v", err)
		http.Error(w, "failed to decode admission review", http.StatusBadRequest)
		return
	}

	if review.Request == nil {
		logrus.Error("Admission review has no request")
		http.Error(w, "missing request", http.StatusBadRequest)
		return
	}

	response := admissionv1.AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: true,
	}

	// Deserialize pod for READING only. Unknown fields (like init container
	// restartPolicy on older k8s.io/api versions) are dropped during
	// deserialization, but that's fine — we only read known fields for
	// decision-making. Mutations use a targeted JSON Patch that never
	// touches fields we didn't explicitly set.
	var pod corev1.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &pod); err != nil {
		logrus.Errorf("Failed to deserialize pod: %v", err)
		response.Result = &metav1.Status{Message: fmt.Sprintf("failed to deserialize pod: %v", err)}
		writeAdmissionResponse(w, &response)
		return
	}

	// Scheduler name is already assigned to a non-default scheduler
	if pod.Spec.SchedulerName != "" && pod.Spec.SchedulerName != "default-scheduler" {
		writeAdmissionResponse(w, &response)
		return
	}

	ctx := r.Context()
	if shouldMutateScheduler(ctx, &pod, review.Request.Namespace, cli, cfg) {
		patch := []jsonPatchOp{{
			Op:    "add",
			Path:  "/spec/schedulerName",
			Value: cfg.schedulerName,
		}}
		patchBytes, err := json.Marshal(patch)
		if err != nil {
			logrus.Errorf("Failed to marshal patch: %v", err)
			writeAdmissionResponse(w, &response)
			return
		}
		patchType := admissionv1.PatchTypeJSONPatch
		response.Patch = patchBytes
		response.PatchType = &patchType
		logrus.Debugf("Mutating pod %s/%s: setting schedulerName=%s",
			review.Request.Namespace, pod.Name, cfg.schedulerName)
	}

	writeAdmissionResponse(w, &response)
}

func writeAdmissionResponse(w http.ResponseWriter, response *admissionv1.AdmissionResponse) {
	review := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: response,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(review); err != nil {
		logrus.Errorf("Failed to encode admission response: %v", err)
	}
}

func run(cli kubernetes.Interface) error {
	cfg := initFlags()

	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", func(w http.ResponseWriter, r *http.Request) {
		handleMutate(w, r, cli, cfg)
	})

	logrus.Infof("Listening on :8080")
	return http.ListenAndServeTLS(":8080", cfg.certFile, cfg.keyFile, mux)
}

func main() {
	cli, err := GetK8sSTDClients()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting kubernetes client: %s", err)
		os.Exit(1)
	}
	err = run(cli)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running app: %s", err)
		os.Exit(1)
	}
}

// GetK8sSTDClients returns the kubernetes clientset using in-cluster config.
func GetK8sSTDClients() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting in-cluster config: %w", err)
	}
	return kubernetes.NewForConfig(config)
}
