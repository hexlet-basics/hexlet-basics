package assetstore_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob"
	"gocloud.dev/blob/memblob"

	"hexletbasics/ent"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/testsupport"
)

var pngHeader = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
}

func newStore(t *testing.T) (*assetstore.Store, *ent.Client, *blob.Bucket) {
	t.Helper()
	db := testsupport.NewClient(t)
	bucket := memblob.OpenBucket(nil)
	t.Cleanup(func() { _ = bucket.Close() })
	return assetstore.New(db, bucket, "https://assets.example.test/"), db, bucket
}

func TestPutAndOpen(t *testing.T) {
	store, db, _ := newStore(t)
	ctx := context.Background()
	data := append(append([]byte{}, pngHeader...), bytes.Repeat([]byte{0}, 600)...)

	got, err := store.Put(ctx, assetstore.Upload{
		Filename: "course-cover.anything",
		Body:     bytes.NewReader(data),
	})
	require.NoError(t, err)

	assert.NotZero(t, got.ID)
	assert.Equal(t, "course-cover.anything", got.Filename)
	assert.Equal(t, "image/png", got.ContentType)
	assert.Equal(t, int64(len(data)), got.ByteSize)
	assert.True(t, strings.HasPrefix(got.URL, "https://assets.example.test/storage/"))
	assert.True(t, strings.HasSuffix(got.URL, ".png"))

	record := db.Attachment.GetX(ctx, got.ID)
	assert.Equal(t, got.Filename, record.Filename)
	assert.Equal(t, got.ContentType, record.ContentType)
	assert.Equal(t, got.ByteSize, record.ByteSize)

	reader, err := store.Open(ctx, record.StorageKey)
	require.NoError(t, err)
	defer func() { _ = reader.Close() }()
	assert.Equal(t, "image/png", reader.ContentType)
	assert.False(t, reader.ModTime.IsZero())
	stored, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, data, stored)

	offset, err := reader.Seek(int64(len(pngHeader)), io.SeekStart)
	require.NoError(t, err)
	assert.Equal(t, int64(len(pngHeader)), offset)
	storedAfterSeek, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, data[len(pngHeader):], storedAfterSeek)
}

func TestPutAcceptsExactSizeLimit(t *testing.T) {
	store, _, _ := newStore(t)
	data := make([]byte, assetstore.MaxUploadBytes)
	copy(data, pngHeader)

	got, err := store.Put(context.Background(), assetstore.Upload{
		Filename: "limit.png",
		Body:     bytes.NewReader(data),
	})
	require.NoError(t, err)
	assert.Equal(t, assetstore.MaxUploadBytes, got.ByteSize)
}

func TestPutRejectsUnsupportedContentWithoutWriting(t *testing.T) {
	store, _, bucket := newStore(t)

	_, err := store.Put(context.Background(), assetstore.Upload{
		Filename: "pretends-to-be.png",
		Body:     strings.NewReader("<svg><script>alert(1)</script></svg>"),
	})

	require.ErrorIs(t, err, assetstore.ErrUnsupportedMediaType)
	assertBucketEmpty(t, bucket)
}

func TestPutCleansUpOversizedBlob(t *testing.T) {
	store, _, bucket := newStore(t)
	data := make([]byte, assetstore.MaxUploadBytes+1)
	copy(data, pngHeader)

	_, err := store.Put(context.Background(), assetstore.Upload{
		Filename: "large.png",
		Body:     bytes.NewReader(data),
	})

	require.ErrorIs(t, err, assetstore.ErrTooLarge)
	assertBucketEmpty(t, bucket)
}

func TestPutCleansUpAfterReaderFailure(t *testing.T) {
	store, _, bucket := newStore(t)
	data := append(append([]byte{}, pngHeader...), bytes.Repeat([]byte{0}, 1016)...)

	_, err := store.Put(context.Background(), assetstore.Upload{
		Filename: "broken.png",
		Body:     &failAfterFirstRead{data: data},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "source read failed")
	assertBucketEmpty(t, bucket)
}

func TestPutCleansUpAfterAttachmentFailure(t *testing.T) {
	store, db, bucket := newStore(t)
	db.Attachment.Use(func(ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(context.Context, ent.Mutation) (ent.Value, error) {
			return nil, errors.New("attachment insert failed")
		})
	})
	data := append(append([]byte{}, pngHeader...), bytes.Repeat([]byte{0}, 600)...)

	_, err := store.Put(context.Background(), assetstore.Upload{
		Filename: "orphan.png",
		Body:     bytes.NewReader(data),
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "attachment insert failed")
	assertBucketEmpty(t, bucket)
}

func TestOpenUnknownKey(t *testing.T) {
	store, _, _ := newStore(t)

	_, err := store.Open(context.Background(), "missing.png")

	require.ErrorIs(t, err, assetstore.ErrNotFound)
}

type failAfterFirstRead struct {
	data []byte
	read bool
}

func (r *failAfterFirstRead) Read(p []byte) (int, error) {
	if r.read {
		return 0, errors.New("source read failed")
	}
	r.read = true
	return copy(p, r.data), nil
}

func assertBucketEmpty(t *testing.T, bucket *blob.Bucket) {
	t.Helper()
	_, err := bucket.List(nil).Next(context.Background())
	require.ErrorIs(t, err, io.EOF)
}
