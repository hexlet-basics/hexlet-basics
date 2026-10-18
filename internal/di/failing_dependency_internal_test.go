package di

import (
	"errors"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
)

var errPoisoned = errors.New("poisoned dependency")

// directDependencies resolves every service of a container once, so the
// injector's graph records who asked for whom, and returns each service's direct
// dependencies. The wiring is read from the graph rather than restated here, so a
// provider added tomorrow is covered without touching this test.
func directDependencies(t *testing.T, newContainer func() *do.RootScope) map[string][]string {
	t.Helper()

	injector := newContainer()

	t.Cleanup(func() { _ = injector.Shutdown() })

	deps := map[string][]string{}

	for _, service := range injector.ListProvidedServices() {
		_, err := do.InvokeNamed[any](injector, service.Service)
		require.NoError(t, err, service.Service)
	}

	for _, service := range injector.ListProvidedServices() {
		explained, ok := do.ExplainNamedService(injector, service.Service)
		require.True(t, ok, service.Service)

		for _, dependency := range explained.Dependencies {
			deps[service.Service] = append(deps[service.Service], dependency.Service)
		}
	}

	return deps
}

// A provider that cannot get one of its dependencies must fail with it, not build
// a half-wired service. Each service is resolved once per dependency with that one
// dependency replaced by a failing provider; none may come back without an error.
func testProvidersReportFailingDependency(t *testing.T, newContainer func() *do.RootScope) {
	t.Helper()
	configureTestEnvironment(t)

	for service, dependencies := range directDependencies(t, newContainer) {
		for _, dependency := range dependencies {
			t.Run(service+" without "+dependency, func(t *testing.T) {
				injector := newContainer()

				t.Cleanup(func() { _ = injector.Shutdown() })
				do.OverrideNamed(injector, dependency, func(do.Injector) (any, error) {
					return nil, errPoisoned
				})

				_, err := do.InvokeNamed[any](injector, service)

				require.Error(t, err)
			})
		}
	}
}

func TestServerProvidersReportFailingDependency(t *testing.T) {
	testProvidersReportFailingDependency(t, newServerContainer)
}

func TestWorkerProvidersReportFailingDependency(t *testing.T) {
	testProvidersReportFailingDependency(t, newWorkerContainer)
}
