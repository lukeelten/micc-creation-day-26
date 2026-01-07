package pkg

import (
	"context"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
)

type Application struct {
	Pb        *pocketbase.PocketBase
	K8sClient *kubernetes.Clientset
}

func NewApplication() (*Application, error) {
	clientset, err := NewKubernetesClient()
	if err != nil {
		return nil, err
	}

	pb := pocketbase.New()
	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {

		go runMetricsServer(pb.RootCmd.Context())

		e.Router.GET("/{path...}", apis.Static(os.DirFS("./public"), true))

		return e.Next()
	})

	return &Application{
		Pb:        pb,
		K8sClient: clientset,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	a.Pb.RootCmd.SetContext(ctx)

	return a.Pb.Start()
}

func runMetricsServer(ctx context.Context) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    ":9090",
		Handler: promhttp.Handler(),
	}

	<-ctx.Done()
	server.Shutdown(context.Background())
}
