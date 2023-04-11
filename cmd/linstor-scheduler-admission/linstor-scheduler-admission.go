package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
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

func run(cli kubernetes.Interface) error {
	logrusLogEntry := logrus.NewEntry(logrus.New())
	logrusLogEntry.Logger.SetLevel(logrus.DebugLevel)
	logger := kwhlogrus.NewLogrus(logrusLogEntry)

	cfg := initFlags()

	// Create mutator.
	mt := kwhmutating.MutatorFunc(func(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
		pod, ok := obj.(*corev1.Pod)
		if !ok {
			return &kwhmutating.MutatorResult{}, nil
		}

		// Scheduler name is already assigned
		if pod.Spec.SchedulerName != "" && pod.Spec.SchedulerName != "default-scheduler" {
			return &kwhmutating.MutatorResult{}, nil
		}

		var pvcNames []string

		// Collect all PVCs attached to Pod
		for _, volume := range pod.Spec.Volumes {
			// Volume has inline CSI driver assigned
			if volume.CSI != nil && volume.CSI.Driver == cfg.driverName {
				pod.Spec.SchedulerName = cfg.schedulerName
				return &kwhmutating.MutatorResult{MutatedObject: pod}, nil
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
			pvc, err := cli.CoreV1().PersistentVolumeClaims(ar.Namespace).Get(ctx, pvcName, metav1.GetOptions{})
			if err != nil && !errors.IsNotFound(err) {
				return &kwhmutating.MutatorResult{}, err
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
			if discoveredProvisioner == "" && pvc.Spec.StorageClassName != nil && *pvc.Spec.StorageClassName != "" {
				sc, err := cli.StorageV1().StorageClasses().Get(ctx, *pvc.Spec.StorageClassName, metav1.GetOptions{})
				if err != nil && !errors.IsNotFound(err) {
					return &kwhmutating.MutatorResult{}, err
				}
				if sc != nil && sc.Provisioner == cfg.driverName {
					discoveredProvisioner = sc.Provisioner
				}
			}
			// Try to gather provisioner name from associated PV
			if discoveredProvisioner == "" && pvc.Spec.VolumeName != "" {
				pv, err := cli.CoreV1().PersistentVolumes().Get(ctx, pvc.Spec.VolumeName, metav1.GetOptions{})
				if err != nil && !errors.IsNotFound(err) {
					return &kwhmutating.MutatorResult{}, err
				}
				if pv != nil && pv.Spec.CSI != nil {
					discoveredProvisioner = pv.Spec.CSI.Driver
				}
			}
			// Overwrite the scheduler name
			if discoveredProvisioner == cfg.driverName {
				pod.Spec.SchedulerName = cfg.schedulerName
				break
			}
		}

		return &kwhmutating.MutatorResult{MutatedObject: pod}, nil
	})

	// Create webhook.
	mcfg := kwhmutating.WebhookConfig{
		ID:      "linstor-scheduler-admission",
		Mutator: mt,
		Logger:  logger,
	}
	wh, err := kwhmutating.NewWebhook(mcfg)
	if err != nil {
		return fmt.Errorf("error creating webhook: %w", err)
	}

	// Get HTTP handler from webhook.
	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{Webhook: wh, Logger: logger})
	if err != nil {
		return fmt.Errorf("error creating webhook handler: %w", err)
	}

	// Serve.
	logger.Infof("Listening on :8080")
	err = http.ListenAndServeTLS(":8080", cfg.certFile, cfg.keyFile, whHandler)
	if err != nil {
		return fmt.Errorf("error serving webhook: %w", err)
	}

	return nil
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

// GetK8sSTDClients returns a all k8s clients.
func GetK8sSTDClients() (kubernetes.Interface, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// Get the client.
	k8sCli, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return k8sCli, nil
}
