package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/kvaps/linstor-scheduler-extender/pkg/consts"
	_ "github.com/kvaps/linstor-scheduler-extender/pkg/driver"
	"github.com/libopenstorage/stork/drivers/volume"
	"github.com/libopenstorage/stork/pkg/extender"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	api_v1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	core_v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var ext *extender.Extender

func main() {
	// Parse empty flags to suppress warnings from the snapshotter which uses
	// glog
	err := flag.CommandLine.Parse([]string{})
	if err != nil {
		log.Warnf("Error parsing flag: %v", err)
	}
	err = flag.Set("logtostderr", "true")
	if err != nil {
		log.Fatalf("Error setting glog flag: %v", err)
	}

	app := cli.NewApp()
	app.Name = consts.Name
	app.Version = consts.Version
	app.Usage = "Linstor scheduler extender for Kubernetes"
	app.Action = run

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Enable verbose logging",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error starting: %v", err)
	}
}

func run(c *cli.Context) {

	log.WithField("version", consts.Version).Infof("starting " + consts.Name)

	verbose := c.Bool("verbose")
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting cluster config: %v", err)
	}

	k8sClient, err := clientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error getting client, %v", err)
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(&core_v1.EventSinkImpl{Interface: k8sClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, api_v1.EventSource{Component: consts.EventComponentName})

	d, err := volume.Get("linstor")
	if err != nil {
		log.Fatalf("Error getting Stork Driver %v: %v", "linstor", err)
	}

	if err = d.Init(nil); err != nil {
		log.Fatalf("Error initializing Stork Driver %v: %v", "linstor", err)
	}

	ext = &extender.Extender{
		Driver:   d,
		Recorder: recorder,
	}

	if err = ext.Start(); err != nil {
		log.Fatalf("Error starting scheduler extender: %v", err)
	}
	// Create operator-sdk manager that will manage all controllers.
	mgr, err := manager.New(config, manager.Options{})
	if err != nil {
		log.Fatalf("Setup controller manager: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()

	go func() {
		for {
			<-signalChan
			log.Printf("Shutdown signal received, exiting...")
			if c.Bool("extender") {
				if err := ext.Stop(); err != nil {
					log.Warnf("Error stopping extender: %v", err)
				}
			}
			if err := d.Stop(); err != nil {
				log.Warnf("Error stopping driver: %v", err)
			}
			ctx.Done()
		}
	}()

	if err := mgr.Start(ctx); err != nil {
		log.Fatalf("Controller manager: %v", err)
	}
	os.Exit(0)
}
