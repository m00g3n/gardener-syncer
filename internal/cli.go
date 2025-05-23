package cli

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-project/gardener-syncer/internal/k8s/client"
	seeker "github.com/kyma-project/gardener-syncer/pkg"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var defaultKcpClientTimeout = time.Second * 10

func Run() error {
	cfg, err := NewConfigFromFlags()
	if err != nil {
		return err
	}
	slog.Info("application started")

	kcpClient, err := client.New(client.Options{
		AdditionalAddToSchema: []func(*runtime.Scheme) error{
			corev1.AddToScheme,
		},
	})

	if err != nil {
		return err
	}

	store := seeker.BuildStoreFn(seeker.StoreOpts{
		Key:     cfg.seedMapKey(),
		Patch:   kcpClient.Patch,
		Get:     kcpClient.Get,
		Convert: seeker.ToConfigMap,
		Timeout: defaultKcpClientTimeout,
	})

	gardenerClient, err := client.New(
		client.Options{
			KubeconfigPath: cfg.Gardener.KubeconfigPath,
			AdditionalAddToSchema: []func(*runtime.Scheme) error{
				v1beta1.AddToScheme,
			},
		},
	)

	if err != nil {
		return err
	}

	gardenerTimeout := mustParseDuration(cfg.Gardener.Timeout)
	fetch := seeker.BuildFetchSeedFn(seeker.FetchSeedsOpts{
		List:    gardenerClient.List,
		Timeout: gardenerTimeout,
	})

	sync := seeker.BuildSyncFn(store, fetch)
	return sync()
}

func mustParseDuration(s string) time.Duration {
	out, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("invalid duration value: %s", s))
	}
	return out
}
