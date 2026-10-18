package apiconv

import (
	"encoding/json"

	"github.com/go-faster/jx"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/audit"
)

// AuditRecord converts one stored Audit Record to its API shape. It is written
// by hand rather than generated because `changes` and `reads` are JSON columns
// with a shape of their own (internal/audit.Change and Reads), and the values
// of a changed field are JSON of any type that ogen carries as raw bytes.
// A record with no changes (a read, a refusal) lists an empty `changes`; one
// with no reads has `reads` null.
func AuditRecord(row *ent.AuditRecord) (api.AuditRecord, error) {
	out := api.AuditRecord{
		ID:             int64(row.ID),
		OccurredAt:     row.OccurredAt,
		ActorUserId:    NilInt32FromPtr(row.ActorUserID),
		Channel:        api.AuditRecordChannel(row.Channel),
		ConnectedAppId: NilInt32FromPtr(row.ConnectedAppID),
		TokenId:        NilStringFromPtr(row.TokenID),
		RequestId:      NilStringFromPtr(row.RequestID),
		Address:        NilStringFromPtr(row.Address),
		UserAgent:      NilStringFromPtr(row.UserAgent),
		Operation:      row.Operation,
		EntityType:     NilStringFromPtr(row.EntityType),
		EntityId:       NilStringFromPtr(row.EntityID),
		Locale:         NilStringFromPtr(row.Locale),
		Outcome:        api.AuditRecordOutcome(row.Outcome),
		Changes:        []api.AuditChange{},
		Reads:          api.NilAuditReads{Null: true},
	}

	if len(row.Changes) > 0 {
		var changes []audit.Change
		if err := json.Unmarshal(row.Changes, &changes); err != nil {
			return api.AuditRecord{}, errs.With("audit_record_id", row.ID).Wrapf(err, "decode audit changes")
		}
		for _, change := range changes {
			item := api.AuditChange{Entity: change.Entity, Op: change.Op}
			if change.ID != "" {
				item.ID = api.NewOptString(change.ID)
			}
			for _, field := range change.Fields {
				item.Fields = append(item.Fields, api.AuditFieldChange{
					Name:   field.Name,
					Old:    jx.Raw(field.Old),
					New:    jx.Raw(field.New),
					Secret: api.NewOptBool(field.Secret),
				})
			}
			out.Changes = append(out.Changes, item)
		}
	}

	if len(row.Reads) > 0 {
		var reads audit.Reads
		if err := json.Unmarshal(row.Reads, &reads); err != nil {
			return api.AuditRecord{}, errs.With("audit_record_id", row.ID).Wrapf(err, "decode audit reads")
		}
		out.Reads = api.NewNilAuditReads(api.AuditReads{Count: Int32FromInt(reads.Count), Ids: reads.IDs})
	}
	return out, nil
}
