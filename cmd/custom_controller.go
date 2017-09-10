package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prateekpandey14/kubecontroller/pkg/controller"
	"github.com/prateekpandey14/kubecontroller/version"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.Printf("custom-controller version %s", version.VERSION)

	kubeconfig := ""
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "kubeconfig file")
	flag.Parse()
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}
	var (
		config *rest.Config
		err    error
	)
	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client: %v", err)
		os.Exit(1)
	}
	client := kubernetes.NewForConfigOrDie(config)

	sharedInformers := informers.NewSharedInformerFactory(client, 10*time.Minute)
	cController := controller.NewCustomController(client, sharedInformers.Core().V1().Secrets(), sharedInformers.Core().V1().Namespaces())

	sharedInformers.Start(nil)
	cController.Run(nil)
}
