package handlers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"hexletbasics/internal/testsupport"
)

// What an Audit Record does when the database fails under the operation it
// belongs to. The harness has one transaction and no way to take the database
// down mid-request, so the storage is broken inside that transaction: a table
// renamed away makes the reads that touch it fail, a trigger that raises makes
// the writes to a table fail. DDL is transactional, so the rollback at cleanup
// undoes it. These are the only storage-failure tests left in the package: the
// guarantee they check (no change without its record) is the audit trail's own.

// breakReads makes every read of the table fail.
func breakReads(t *testing.T, h *testsupport.Harness, table string) {
	t.Helper()
	exec(t, h, fmt.Sprintf(`ALTER TABLE %q RENAME TO %q`, table, table+"_gone"))
}

// breakWrites makes every INSERT, UPDATE or DELETE (event) on the table fail,
// while its reads keep working.
func breakWrites(t *testing.T, h *testsupport.Harness, table, event string) {
	t.Helper()
	exec(t, h, `CREATE OR REPLACE FUNCTION storage_down() RETURNS trigger LANGUAGE plpgsql
		AS $$ BEGIN RAISE EXCEPTION 'storage is down'; END $$`)
	exec(t, h, fmt.Sprintf(
		`CREATE TRIGGER storage_down BEFORE %s ON %q FOR EACH ROW EXECUTE FUNCTION storage_down()`,
		event, table))
}

func exec(t *testing.T, h *testsupport.Harness, statement string) {
	t.Helper()
	_, err := h.Tx.ExecContext(t.Context(), statement)
	require.NoError(t, err)
}
