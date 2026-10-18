package mcpserver_test

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/mcpserver"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/testsupport"
)

// tinyPNG is a minimal (1x1) PNG, real image bytes for the store's sniffing.
func tinyPNG() []byte {
	return []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
		0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}
}

func upload(fileName string, content []byte) map[string]any {
	return map[string]any{"locale": "en", "body": map[string]any{
		"fileName": fileName, "contentType": "image/png",
		"contentBase64": base64.StdEncoding.EncodeToString(content),
	}}
}

func attachments(t *testing.T, o *testsupport.OAuth) int {
	t.Helper()

	return o.DB.Attachment.Query().CountX(t.Context())
}

func TestUploadToolStoresTheAttachmentAndRecordsItOnTheMCPChannel(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor(), "mcp.write")
	session, token := admin(t, o, bob())
	before := attachments(t, o)

	result := call(t, session, "adminUploadAttachment", upload("cover.png", tinyPNG()))

	require.False(t, result.IsError, "%+v", result.Content)
	out := as[map[string]any](t, result.StructuredContent)
	assert.Equal(t, "cover.png", out["filename"])
	assert.Equal(t, "image/png", out["contentType"])
	assert.EqualValues(t, len(tinyPNG()), out["byteSize"])
	assert.NotEmpty(t, out["url"])
	assert.Equal(t, before+1, attachments(t, o))

	records := o.Records("AdminUploadAttachment")
	require.Len(t, records, 1)
	assert.Equal(t, audit.OutcomeSucceeded, audit.Outcome(records[0].Outcome))
	assert.Equal(t, "mcp", records[0].Channel)
	assert.Equal(t, bob(), *records[0].ActorUserID)
	assert.Equal(t, token.ID, *records[0].TokenID)
	assert.Equal(t, token.ConnectedAppID, *records[0].ConnectedAppID)
	assert.Equal(t, strconv.Itoa(int(as[float64](t, out["id"]))), *records[0].EntityID)
}

func TestUploadToolOfferedWithTheContractsPartAndLimits(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor(), "mcp.write")
	session, _ := admin(t, o, bob())
	tools, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	for _, tool := range tools.Tools {
		if tool.Name != "adminUploadAttachment" {
			continue
		}

		assert.False(t, tool.Annotations.ReadOnlyHint)
		assert.Contains(t, schemaJSON(t, tool.InputSchema), "contentBase64")
		assert.Contains(t, schemaJSON(t, tool.OutputSchema), "byteSize")

		return
	}

	t.Fatal("adminUploadAttachment is not offered")
}

func TestUploadIsRefusedAndRecordedAsForbiddenWithoutWhatTheWebsiteNeeds(t *testing.T) {
	cases := map[string]func(t *testing.T, o *testsupport.OAuth) *mcp.ClientSession{
		"no mcp.write": func(t *testing.T, o *testsupport.OAuth) *mcp.ClientSession {
			t.Helper()
			o.Grant(editor(), "mcp.read")
			session, _ := admin(t, o, bob())

			return session
		},
		"no admin scope": func(t *testing.T, o *testsupport.OAuth) *mcp.ClientSession {
			t.Helper()
			o.Grant(editor(), "mcp.write")

			return connect(t, o, bob(), oauthserver.ScopeCatalogRead)
		},
		"not staff": func(t *testing.T, o *testsupport.OAuth) *mcp.ClientSession {
			t.Helper()
			session, _ := admin(t, o, alice())

			return session
		},
	}
	for name, setup := range cases {
		t.Run(name, func(t *testing.T) {
			o := testsupport.NewOAuth(t)
			session := setup(t, o)
			before := attachments(t, o)

			assertRefused(t, call(t, session, "adminUploadAttachment", upload("cover.png", tinyPNG())))

			assert.Equal(t, before, attachments(t, o), "nothing is stored")
			records := o.Records("AdminUploadAttachment")
			require.Len(t, records, 1)
			assert.Equal(t, "forbidden", records[0].Outcome)
			assert.Equal(t, "mcp", records[0].Channel)
		})
	}
}

func TestUploadOfTheLargestAllowedFileFitsAndOneByteMoreIsRefusedBeforeDecoding(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor(), "mcp.write")
	session, _ := admin(t, o, bob())
	before := attachments(t, o)

	largest := append(append([]byte{}, tinyPNG()...), make([]byte, int(assetstore.MaxUploadBytes)-len(tinyPNG()))...)
	ok := call(t, session, "adminUploadAttachment", upload("big.png", largest))
	require.False(t, ok.IsError, "%+v", ok.Content)
	assert.Equal(t, before+1, attachments(t, o))

	tooBig := append(append([]byte{}, largest...), 0)
	result := call(t, session, "adminUploadAttachment", upload("huge.png", tooBig))
	require.True(t, result.IsError)
	assert.Contains(t, textOf(t, result), "too large")
	assert.Equal(t, before+1, attachments(t, o), "the oversized file is not stored")
}

func TestUploadRejectsWhatTheWebsiteRejects(t *testing.T) {
	o := testsupport.NewOAuth(t)
	o.Grant(editor(), "mcp.write")
	session, _ := admin(t, o, bob())
	before := attachments(t, o)

	t.Run("not an image", func(t *testing.T) {
		result := call(t, session, "adminUploadAttachment", upload("notes.txt", []byte("plain text")))
		require.True(t, result.IsError)
		assert.Contains(t, textOf(t, result), "HTTP 422")
	})
	t.Run("invalid base64", func(t *testing.T) {
		args := upload("cover.png", tinyPNG())
		as[map[string]any](t, args["body"])["contentBase64"] = "not*base64!"
		result := call(t, session, "adminUploadAttachment", args)
		require.True(t, result.IsError)
		assert.Contains(t, textOf(t, result), "not valid base64")
	})
	t.Run("missing file name", func(t *testing.T) {
		result := call(t, session, "adminUploadAttachment", upload("", tinyPNG()))
		require.True(t, result.IsError)
		assert.Contains(t, textOf(t, result), "fileName")
	})
	assert.Equal(t, before, attachments(t, o))
}

func TestUploadSpendsTheWriteBudget(t *testing.T) {
	o := testsupport.NewOAuth(t, testsupport.WithMCPLimits(mcpserver.Limits{User: mcpserver.Rates{Write: 1}}))
	o.Grant(editor(), "mcp.write")
	session := connect(t, o, bob(), oauthserver.ScopeAdmin)

	require.False(t, callTool(t, session, "adminUploadAttachment", upload("a.png", tinyPNG())).IsError)
	rejected(t, callTool(t, session, "adminUploadAttachment", upload("b.png", tinyPNG())))
	assert.Len(t, o.Records("AdminUploadAttachment"), 1)
}

func schemaJSON(t *testing.T, schema any) string {
	t.Helper()

	raw, err := json.Marshal(schema)
	require.NoError(t, err)

	return string(raw)
}
