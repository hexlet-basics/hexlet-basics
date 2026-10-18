package mailer

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/config"
)

// postboxSent is what the fake Postbox endpoint was asked to deliver.
type postboxSent struct {
	From    string
	To      []string
	Subject string
	HTML    string
	Path    string
}

// fakePostbox stands in for the provider's SES v2 API at the HTTP boundary, so
// Postbox runs for real — request building, signing, response handling — and the
// test reads what the provider would have received.
func fakePostbox(t *testing.T, status int) (config.MailConfig, func() []postboxSent) {
	t.Helper()
	var (
		mu   sync.Mutex
		sent []postboxSent
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		var request struct {
			FromEmailAddress string `json:"FromEmailAddress"`
			Destination      struct {
				ToAddresses []string `json:"ToAddresses"`
			} `json:"Destination"`
			Content struct {
				Simple struct {
					Subject struct {
						Data string `json:"Data"`
					} `json:"Subject"`
					Body struct {
						Html struct {
							Data string `json:"Data"`
						} `json:"Html"`
					} `json:"Body"`
				} `json:"Simple"`
			} `json:"Content"`
		}
		assert.NoError(t, json.Unmarshal(body, &request))
		mu.Lock()
		sent = append(sent, postboxSent{
			From:    request.FromEmailAddress,
			To:      request.Destination.ToAddresses,
			Subject: request.Content.Simple.Subject.Data,
			HTML:    request.Content.Simple.Body.Html.Data,
			Path:    r.URL.Path,
		})
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"MessageId":"m-1","message":"refused"}`))
	}))
	t.Cleanup(server.Close)

	cfg := config.MailConfig{
		From:                   "support@hexlet.io",
		PostboxEndpoint:        server.URL,
		PostboxRegion:          "ru-central1",
		PostboxAccessKeyID:     "key",
		PostboxSecretAccessKey: "secret",
	}
	return cfg, func() []postboxSent {
		mu.Lock()
		defer mu.Unlock()
		return append([]postboxSent(nil), sent...)
	}
}

func TestPostboxDeliversTheMessageAsASimpleHTMLEmail(t *testing.T) {
	t.Parallel()
	cfg, delivered := fakePostbox(t, http.StatusOK)

	err := New(cfg, slog.New(slog.DiscardHandler)).Send(t.Context(), Message{
		To: "learner@example.com", Subject: "Welcome", HTML: "<p>Hello</p>",
	})

	require.NoError(t, err)
	require.Len(t, delivered(), 1)
	assert.Equal(t, postboxSent{
		From:    "support@hexlet.io",
		To:      []string{"learner@example.com"},
		Subject: "Welcome",
		HTML:    "<p>Hello</p>",
		Path:    "/v2/email/outbound-emails",
	}, delivered()[0])
}

func TestPostboxReportsAMessageTheProviderRefused(t *testing.T) {
	t.Parallel()
	cfg, _ := fakePostbox(t, http.StatusBadRequest)

	err := NewPostbox(cfg).Send(t.Context(), Message{To: "learner@example.com", Subject: "Welcome"})

	require.Error(t, err)
	assert.ErrorContains(t, err, "send email via postbox")
}

func TestNewChoosesPostboxOnlyWhenItsAccessKeyIsSet(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.DiscardHandler)

	assert.IsType(t, &Postbox{}, New(config.MailConfig{PostboxAccessKeyID: "key"}, logger))
	assert.IsType(t, &Log{}, New(config.MailConfig{}, logger))
}

func TestLogAdapterWritesTheMessageAndNeverFails(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer

	err := NewLog(slog.New(slog.NewTextHandler(&out, nil))).Send(t.Context(), Message{
		To: "learner@example.com", Subject: "Reset your password", HTML: "<a href=\"/reset\">link</a>",
	})

	require.NoError(t, err)
	assert.Contains(t, out.String(), "learner@example.com")
	assert.Contains(t, out.String(), "Reset your password")
}

func TestRecorderKeepsMessagesInOrderAndHandsOutACopy(t *testing.T) {
	t.Parallel()
	recorder := &Recorder{}
	require.NoError(t, recorder.Send(t.Context(), Message{To: "a@example.com"}))
	require.NoError(t, recorder.Send(t.Context(), Message{To: "b@example.com"}))

	sent := recorder.Sent()
	sent[0].To = "changed@example.com"

	assert.Equal(t, []string{"a@example.com", "b@example.com"},
		[]string{recorder.Sent()[0].To, recorder.Sent()[1].To})
}
