package di

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/health"
)

func probe(t *testing.T, handler http.Handler, path string) int {
	t.Helper()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil))
	return rec.Code
}

func builtHandler(t *testing.T) (do.Injector, http.Handler) {
	t.Helper()
	configureTestEnvironment(t)
	injector := newServerContainer()
	t.Cleanup(func() { _ = injector.Shutdown() })
	handler, err := do.InvokeNamed[http.Handler](injector, httpHandlerServiceName)
	require.NoError(t, err)
	return injector, handler
}

// failDatabase makes every query the client runs fail, as a database that went
// away would, without touching DDL.
func failDatabase(t *testing.T, injector do.Injector) {
	t.Helper()
	client, err := do.Invoke[*ent.Client](injector)
	require.NoError(t, err)
	client.Intercept(ent.InterceptFunc(func(ent.Querier) ent.Querier {
		return ent.QuerierFunc(func(context.Context, ent.Query) (ent.Value, error) {
			return nil, errors.New("database is down")
		})
	}))
}

func TestHealthzAndReadyzAreOKWhileTheProcessServes(t *testing.T) {
	_, handler := builtHandler(t)

	assert.Equal(t, http.StatusOK, probe(t, handler, "/healthz"))
	assert.Equal(t, http.StatusOK, probe(t, handler, "/readyz"))
}

func TestReadyzIs503WhenTheDatabaseFailsButHealthzStaysOK(t *testing.T) {
	injector, handler := builtHandler(t)
	failDatabase(t, injector)

	assert.Equal(t, http.StatusOK, probe(t, handler, "/healthz"))
	assert.Equal(t, http.StatusServiceUnavailable, probe(t, handler, "/readyz"))
}

func TestReadyzIs503AfterTheProcessStartsDraining(t *testing.T) {
	injector, handler := builtHandler(t)
	probes, err := do.Invoke[*health.Probes](injector)
	require.NoError(t, err)

	probes.Drain()

	assert.Equal(t, http.StatusServiceUnavailable, probe(t, handler, "/readyz"))
	assert.Equal(t, http.StatusOK, probe(t, handler, "/healthz"))
}

func TestWorkerHealthServerIsReadyOnlyAfterStartAndWhileTheDatabaseWorks(t *testing.T) {
	configureWorkerEnvironment(t)
	injector := newWorkerContainer()
	t.Cleanup(func() { _ = injector.Shutdown() })
	server, err := do.InvokeNamed[*http.Server](injector, healthServerServiceName)
	require.NoError(t, err)
	probes, err := do.Invoke[*health.Probes](injector)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, probe(t, server.Handler, "/healthz"))
	assert.Equal(t, http.StatusServiceUnavailable, probe(t, server.Handler, "/readyz"), "before start")

	probes.Started()
	assert.Equal(t, http.StatusOK, probe(t, server.Handler, "/readyz"))

	failDatabase(t, injector)
	assert.Equal(t, http.StatusServiceUnavailable, probe(t, server.Handler, "/readyz"), "database down")
	assert.Equal(t, http.StatusOK, probe(t, server.Handler, "/healthz"))
}
