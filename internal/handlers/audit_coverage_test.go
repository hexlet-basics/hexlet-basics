package handlers_test

import (
	"errors"
	"testing"

	"github.com/ogen-go/ogen/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/audit"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/permissioncatalog"
	"hexletbasics/internal/testsupport"
)

// Every admin write operation of the contract leaves an Audit Record, whatever
// marks it (a write Permission, an administrators-only or any-staff marker, or
// only the admin session scheme). The operations are enumerated from the
// contract, so adding one to the TypeSpec adds it here with no edit, and one the
// middleware does not reach fails this test.
func TestEveryAdminWriteOperationOfTheContractIsRecorded(t *testing.T) {
	operations := permissioncatalog.AdminWrites()
	require.NotEmpty(t, operations)

	for _, operation := range operations {
		t.Run(operation, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			audited := auditMiddleware(t, h)

			_, err := audited(middleware.Request{Context: t.Context(), OperationName: operation}, succeeds)
			require.NoError(t, err)

			rows := records(t, h, operation)
			require.Len(t, rows, 1, "%s is an admin write and must be recorded", operation)
			assert.Equal(t, "succeeded", rows[0].Outcome)
		})
	}
}

// A failed call records outcome failed, in a transaction of its own so the
// record survives the rollback of the change. It holds no changes, because the
// changes did not happen and may name values the caller never saw.
func TestEveryAdminWriteOperationRecordsAFailure(t *testing.T) {
	for _, operation := range permissioncatalog.AdminWrites() {
		t.Run(operation, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			audited := auditMiddleware(t, h)

			_, err := audited(middleware.Request{Context: t.Context(), OperationName: operation}, func(middleware.Request) (middleware.Response, error) {
				return middleware.Response{}, errors.New("boom: secret value 4242")
			})
			require.Error(t, err)

			rows := records(t, h, operation)
			require.Len(t, rows, 1)
			assert.Equal(t, string(audit.OutcomeFailed), rows[0].Outcome)
			assert.Empty(t, rows[0].Changes)
			assert.Nil(t, rows[0].EntityID)
		})
	}
}

func succeeds(middleware.Request) (middleware.Response, error) { return middleware.Response{}, nil }

func auditMiddleware(t *testing.T, h *testsupport.Harness) middleware.Middleware {
	t.Helper()
	return handlers.NewAuditMiddleware(h.Transactor, testsupport.NewTranslator(t), testsupport.NewAPIErrorHandler(t, testsupport.NewTranslator(t)))
}
