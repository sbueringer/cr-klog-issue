package main

import (
	"context"
	"flag"
	_ "net/http/pprof"
	"os"

	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	// +kubebuilder:scaffold:imports
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	_ "k8s.io/component-base/logs/json/register"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/cluster-api/version"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme     = runtime.NewScheme()
	setupLog   = ctrl.Log.WithName("setup")
	logOptions = logs.NewOptions()
)

func init() {
	_ = corev1.AddToScheme(scheme)
}

// InitFlags initializes the flags.
func InitFlags(fs *pflag.FlagSet) {
	logs.AddFlags(fs, logs.SkipLoggingConfigurationFlags())
	logOptions.AddFlags(fs)
}

func main() {
	InitFlags(pflag.CommandLine)
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := logOptions.ValidateAndApply(); err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// The JSON log format requires the Klog format in klog, otherwise log lines
	// are serialized twice, e.g.:
	// { ... "msg":"controller/cluster \"msg\"=\"Starting workers\"\n"}
	if logOptions.Config.Format == logs.JSONLogFormat {
		ctrl.SetLogger(klogr.NewWithOptions(klogr.WithFormat(klogr.FormatKlog)))
	} else {
		ctrl.SetLogger(klogr.New())
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Setup the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	err = ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(&reconciler{})
	if err != nil {
		setupLog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager", "version", version.Get().String())
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

type reconciler struct{}

func (r *reconciler) Reconcile(ctx context.Context, _ reconcile.Request) (reconcile.Result, error) {

	log := ctrl.LoggerFrom(ctx)

	log.V(0).Info("Test log", "pod", klog.KRef("default", "pod-1"))
	log.V(0).Info("Test log with .String()", "pod", klog.KRef("default", "pod-1").String())

	//log.V(0).Info("Verbosity 0")
	//log.V(1).Info("Verbosity 1")
	//log.V(2).Info("Verbosity 2")
	//log.V(3).Info("Verbosity 3")
	//log.V(4).Info("Verbosity 4")

	return reconcile.Result{}, nil
}
