# Email via Yandex Postbox over its SES-compatible API

Transactional email is sent through **Yandex Cloud Postbox** using its
**AWS SES-compatible API**, called with `aws-sdk-go-v2/service/ses` pointed at
the Postbox endpoint. Despite the AWS SDK dependency, we are NOT on AWS — Postbox
merely emulates the SES API. Templates use Go `html/template`.

Postbox exposes both an SES-compatible API and an SMTP gateway; we pick the SES
API (native, richer, delivery visibility) over SMTP. Legacy sent via plain SMTP
with ERB templates — we deliberately move to the API.

## Consequences

- **DKIM/SPF/DMARC are handled by Postbox**, so no app-side signing — `go-msgauth`
  is dropped.
- Delivery/bounce visibility is via Yandex Cloud Logging/Monitoring, not AWS SNS
  — `robbiet480/go.sns` does not belong to the email path.
- `go-mail` (SMTP) is dropped; it returns only if we later need to build complex
  MIME (attachments).
- Nicer responsive templates (`gomjml`) and a text/plain part (`k3a/html2text`)
  are deferred to a post-parity "improve emails" step, not in the initial set.
