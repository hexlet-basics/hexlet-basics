// Package assetstore owns the complete lifecycle of uploaded assets.
package assetstore

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"

	"hexletbasics/ent"
	"hexletbasics/internal/ids"
)

// MaxUploadBytes is the largest asset accepted by Put.
const MaxUploadBytes int64 = 15 << 20

const sniffBytes = 512

var (
	// ErrUnsupportedMediaType means the uploaded bytes are not an allowed raster image.
	ErrUnsupportedMediaType = errors.New("unsupported asset media type")
	// ErrTooLarge means the uploaded bytes exceed MaxUploadBytes.
	ErrTooLarge = errors.New("asset is too large")
	// ErrNotFound means the requested storage key does not exist.
	ErrNotFound = errors.New("asset not found")
)

var extensionsByContentType = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// Upload is the minimal input Put needs. Content type and byte size are derived
// from Body rather than trusted from caller-supplied metadata.
type Upload struct {
	Filename string
	Body     io.Reader
}

// Attachment describes a successfully stored blob and its database record.
type Attachment struct {
	ID          int
	URL         string
	Filename    string
	ContentType string
	ByteSize    int64
}

// Reader exposes stored bytes and the metadata needed by the HTTP read path
// without leaking the blob backend or its error model to callers.
type Reader struct {
	io.ReadSeekCloser
	ContentType string
	ModTime     time.Time
}

// Store coordinates blob storage and Attachment persistence as one operation.
type Store struct {
	db        *ent.Client
	bucket    *blob.Bucket
	publicURL string
}

// New constructs an asset store. publicURL is the origin serving /storage.
func New(db *ent.Client, bucket *blob.Bucket, publicURL string) *Store {
	return &Store{
		db:        db,
		bucket:    bucket,
		publicURL: strings.TrimRight(publicURL, "/"),
	}
}

// Put validates and streams an image to blob storage, then persists its
// Attachment. Any failure after opening the blob writer compensates by deleting
// the key, so callers never need to coordinate partial storage state.
func (s *Store) Put(ctx context.Context, input Upload) (Attachment, error) {
	src := bufio.NewReader(input.Body)
	header, err := src.Peek(sniffBytes)
	if err != nil && !errors.Is(err, io.EOF) {
		return Attachment{}, fmt.Errorf("inspect asset %q: %w", input.Filename, err)
	}

	contentType := http.DetectContentType(header)
	extension, ok := extensionsByContentType[contentType]
	if !ok {
		return Attachment{}, fmt.Errorf("%w: %s", ErrUnsupportedMediaType, contentType)
	}

	key := ids.New() + extension
	writer, err := s.bucket.NewWriter(ctx, key, &blob.WriterOptions{ContentType: contentType})
	if err != nil {
		return Attachment{}, fmt.Errorf("open blob writer for %q: %w", key, err)
	}

	limited := &io.LimitedReader{R: src, N: MaxUploadBytes + 1}
	written, copyErr := io.Copy(writer, limited)
	closeErr := writer.Close()
	if copyErr != nil || closeErr != nil {
		cause := errors.Join(copyErr, closeErr)
		return Attachment{}, s.compensate(key, fmt.Errorf("write blob %q: %w", key, cause))
	}
	if written > MaxUploadBytes {
		return Attachment{}, s.compensate(key, ErrTooLarge)
	}

	record, err := s.db.Attachment.Create().
		SetStorageKey(key).
		SetFilename(input.Filename).
		SetContentType(contentType).
		SetByteSize(written).
		Save(ctx)
	if err != nil {
		return Attachment{}, s.compensate(key, fmt.Errorf("record attachment %q: %w", key, err))
	}

	return Attachment{
		ID:          record.ID,
		URL:         s.publicURL + "/storage/" + url.PathEscape(key),
		Filename:    record.Filename,
		ContentType: record.ContentType,
		ByteSize:    record.ByteSize,
	}, nil
}

// Open returns a stored asset while translating backend-specific not-found
// errors into the module's stable error vocabulary.
func (s *Store) Open(ctx context.Context, key string) (*Reader, error) {
	reader, err := s.bucket.NewReader(ctx, key, nil)
	if err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			return nil, fmt.Errorf("%w: %s", ErrNotFound, key)
		}
		return nil, fmt.Errorf("open asset %q: %w", key, err)
	}
	return &Reader{
		ReadSeekCloser: reader,
		ContentType:    reader.ContentType(),
		ModTime:        reader.ModTime(),
	}, nil
}

func (s *Store) compensate(key string, cause error) error {
	cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.bucket.Delete(cleanupCtx, key); err != nil && gcerrors.Code(err) != gcerrors.NotFound {
		return errors.Join(cause, fmt.Errorf("delete partial blob %q: %w", key, err))
	}
	return cause
}
