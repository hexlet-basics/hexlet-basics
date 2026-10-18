package di

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/samber/do/v2"
	"gocloud.dev/blob"

	"hexletbasics/ent"
)

// Closer services release the vendor resources that do not implement a do
// shutdown method themselves (ADR-0023). Each is a tiny service holding one
// resource, so the rest of the code keeps seeing the vendor type
// (*ent.Client, *blob.Bucket, *sentry.Client) and never a wrapper.
//
// do shuts a service down only after everything that depends on it has been
// shut down, so a closer that resolves another closer inside its provider
// runs before it. That is how Sentry is flushed last: the database and the
// bucket closers depend on the Sentry closer, so an error they log while
// closing still reaches Sentry.
//
// *http.Server and *otelconf.SDK need no closer: both already have a
// Shutdown(ctx) error method, which do calls on its own.

type databaseCloser struct{ client *ent.Client }

func (c *databaseCloser) Shutdown() error {
	err := c.client.Close()
	if err != nil {
		return errs.Wrapf(err, "closing database")
	}

	return nil
}

// HealthCheck is the readiness probe's database check (ADR-0023): one cheap
// query through the pool, which fails when the database is unreachable.
func (c *databaseCloser) HealthCheck(ctx context.Context) error {
	_, err := c.client.User.Query().Exist(ctx)
	if err != nil {
		return errs.Wrapf(err, "pinging database")
	}

	return nil
}

type bucketCloser struct{ bucket *blob.Bucket }

func (c *bucketCloser) Shutdown() error {
	err := c.bucket.Close()
	if err != nil {
		return errs.Wrapf(err, "closing blob bucket")
	}

	return nil
}

type sentryCloser struct{ client *sentry.Client }

// errSentryFlushTimeout is what the closer answers when Sentry's queue was not
// delivered before the shutdown deadline.
var errSentryFlushTimeout = errors.New("flushing Sentry events timed out")

func (c *sentryCloser) Shutdown(ctx context.Context) error {
	if !c.client.FlushWithContext(ctx) {
		return errSentryFlushTimeout
	}

	return nil
}

// closersPackage registers the closers shared by both process graphs.
var closersPackage = do.Package(
	do.Lazy[*sentryCloser](func(i do.Injector) (*sentryCloser, error) {
		client, err := do.Invoke[*sentry.Client](i)
		if err != nil {
			return nil, err
		}

		return &sentryCloser{client: client}, nil
	}),
	do.Lazy[*databaseCloser](func(i do.Injector) (*databaseCloser, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}

		_, err = do.Invoke[*sentryCloser](i)
		if err != nil {
			return nil, err
		}

		return &databaseCloser{client: client}, nil
	}),
	do.Lazy[*bucketCloser](func(i do.Injector) (*bucketCloser, error) {
		bucket, err := do.Invoke[*blob.Bucket](i)
		if err != nil {
			return nil, err
		}

		_, err = do.Invoke[*sentryCloser](i)
		if err != nil {
			return nil, err
		}

		return &bucketCloser{bucket: bucket}, nil
	}),
)

// registerClosers makes every closer known to the injector, which only shuts
// down services it has built.
func registerClosers(injector do.Injector) error {
	_, err := do.Invoke[*databaseCloser](injector)
	if err != nil {
		return err
	}

	_, err = do.Invoke[*bucketCloser](injector)
	if err != nil {
		return err
	}

	return nil
}

// shutdown releases everything the injector built, after the process's actors
// have stopped. The report names every resource that failed; the caller logs it.
func shutdown(ctx context.Context, injector do.Injector) error {
	report := injector.ShutdownWithContext(ctx)
	if report == nil || report.Succeed {
		return nil
	}

	return report
}
