# Asset storage via gocloud.dev/blob (ActiveStorage replacement)

File/image assets (lesson theory images, avatars) are stored via
**`gocloud.dev/blob`** — a backend-agnostic `blob.Bucket` API — with `s3blob`
in production and `fileblob` (local disk) in development, selected by URL. This
mirrors the ActiveStorage per-environment service split we had, and uses
`aws-sdk-go-v2` under the hood for S3.

`gocloud.dev/blob` is chosen over a bare `aws-sdk-go-v2/service/s3` client for
the dev-local backend swap, and over MinIO/minio-go, which entered maintenance
in Dec 2025 and was archived Feb 2026.

## Consequences

- ActiveStorage is three things — storage + DB attachment records + image
  variants — and no single Go library replaces all three. We compose:
  - **storage** → `gocloud.dev/blob`;
  - **attachment records** (which model owns which blob) → our own ent table;
  - **image variants/resize** → deferred; add an image library
    (`disintegration/imaging`, or libvips `govips`/`bimg`) only if variants turn
    out to be needed (legacy stores theory images as-is).
