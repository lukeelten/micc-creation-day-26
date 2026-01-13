package pkg

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	_ "github.com/lukeelten/micc-creation-day-26/backend/migrations"
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/controller"
	"github.com/lukeelten/micc-creation-day-26/backend/pkg/utils"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
)

type Application struct {
	Pb        *pocketbase.PocketBase
	K8sClient *kubernetes.Clientset
}

func NewApplication() (*Application, error) {
	pb := pocketbase.New()
	kubeconfigFlag := pb.RootCmd.PersistentFlags().String("kubeconfig", "", "(optional) Set manual path to kubeconfig")

	// Start metrics server
	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.App.Logger().Info("Starting metrics server on :9090")
		go runMetricsServer(pb.RootCmd.Context())
		return e.Next()
	})

	// Serve static files
	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		if _, err := os.Stat("./public"); err == nil {
			e.App.Logger().Info("Serving static files from ./public")
			e.Router.GET("/{path...}", apis.Static(os.DirFS("./public"), true))
		} else {
			e.App.Logger().Info("No ./public directory found, skipping static file serving")
		}

		return e.Next()
	})

	// Create superuser if requested via env vars
	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		email, hasEmail := os.LookupEnv("SUPERUSER_EMAIL")
		password, hasPassword := os.LookupEnv("SUPERUSER_PASSWORD")
		if !hasEmail || !hasPassword {
			return e.Next() // skip creating superuser if env vars are not set
		}

		superusers, err := e.App.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}

		record := core.NewRecord(superusers)

		// note: the values can be eventually loaded via os.Getenv(key)
		// or from a special local config file
		record.Set("email", email)
		record.Set("password", password)

		err = e.App.Save(record)
		if err != nil {
			return err
		}

		return e.Next()
	})

	// Bootstrap kubernetes
	pb.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		clientset, err := NewKubernetesClient(e.App.Logger(), kubeconfigFlag)
		if err != nil {
			e.App.Logger().Error("Failed to create kubernetes client", "error", err)
			return err
		}

		pb.Store().Set(utils.StoreClient, clientset)

		return e.Next()
	})

	// Run controller
	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		runController, err := controller.NewRunController(pb)
		if err != nil {
			e.App.Logger().Error("Failed to create RunController", "error", err)
			return err
		}

		if err := runController.Start(); err != nil {
			e.App.Logger().Error("RunController exited with error", "error", err)
			return err
		}

		return e.Next()
	})

	err := SetupApi(pb)
	if err != nil {
		return nil, err
	}

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(pb, pb.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	return &Application{
		Pb: pb,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	a.Pb.RootCmd.SetContext(ctx)
	a.Pb.Store().Set(utils.StoreContext, ctx)

	return a.Pb.Start()
}

func runMetricsServer(ctx context.Context) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	server := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Default().Error("Error during metrics server shutdown", "err", err)
	}
}
