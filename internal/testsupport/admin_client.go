package testsupport

import "hexletbasics/internal/api"

// Client is the generated client over the harness. Every operation answers its
// success type or an error (a non-2xx is a *api.ProblemDetailsStatusCode, ADR-0028),
// so a handler test reads the result directly and an error test inspects the
// returned error.
type Client struct {
	*api.Client
}
