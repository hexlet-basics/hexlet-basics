package mcpserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net/textproto"

	"hexletbasics/internal/assetstore"
)

// A tool takes JSON arguments, which cannot carry a multipart form as it is.
// For an operation whose request body is one binary part (the attachment
// upload), the `body` argument is an object with the file's name, an optional
// content type and the content as base64, and buildRequest turns it into the
// multipart request the website sends. The request then goes through the same
// generated server as any other call (ADR-0021 addendum, #856), so the
// Permission and capability check, validation, the size and type rules of
// assetstore and the Audit Record are the website's, not a second upload path.

// uploadBody is the `body` argument of a multipart tool.
type uploadBody struct {
	FileName      string `json:"fileName"`      //nolint:tagliatelle // an MCP tool argument is named as in the contract, camelCase
	ContentType   string `json:"contentType"`   //nolint:tagliatelle // an MCP tool argument is named as in the contract, camelCase
	ContentBase64 string `json:"contentBase64"` //nolint:tagliatelle // an MCP tool argument is named as in the contract, camelCase
}

// maxUploadBase64 is the longest base64 text of a file the website accepts. It
// is checked on the string, before it is decoded, so an oversized argument never
// costs the memory of its decoded copy.
func maxUploadBase64() int {
	return base64.StdEncoding.EncodedLen(int(assetstore.MaxUploadBytes))
}

// requestBodyLimit is how large a `/mcp` request may be: the SDK's 4 MiB
// default, raised for the upload tool, whose base64 argument for the largest
// allowed file (15 MiB, 20 MiB encoded) would otherwise never arrive. The extra
// room covers the JSON envelope and the other arguments. The limit applies to
// callers who hold a valid token only: the bearer check runs before the body is
// read.
func requestBodyLimit(hasUpload bool) int64 {
	const sdkDefault = 4 << 20
	if !hasUpload {
		return sdkDefault
	}

	return int64(maxUploadBase64()) + 1<<20
}

// uploadSchema describes the `body` argument of a multipart tool for the agent.
func uploadSchema(file *multipartFile) map[string]any {
	return map[string]any{
		"type":                 "object",
		"description":          fmt.Sprintf("The file for the %q part, as base64 text. At most %d bytes once decoded; the content type is sniffed from the bytes and only raster images are accepted, as on the website.", file.Field, assetstore.MaxUploadBytes),
		"additionalProperties": false,
		"required":             []string{"fileName", "contentBase64"},
		"properties": map[string]any{
			"fileName":      map[string]any{"type": "string", "minLength": 1, "description": "The file's name, with its extension."},
			"contentType":   map[string]any{"type": "string", "description": "Optional MIME type of the file. Advisory: the stored type is derived from the bytes."},
			"contentBase64": map[string]any{"type": "string", "minLength": 1, "maxLength": maxUploadBase64(), "description": "The file's bytes, standard base64."},
		},
	}
}

// bytesPerMiB turns the upload limit into the MiB the message names.
const bytesPerMiB = 1 << 20

// multipartBody builds the form body and its Content-Type from the tool's
// `body` argument. A problem is the caller's mistake, reported as a tool error;
// whether the bytes are an allowed image is for the generated server to answer,
// as it does over HTTP.
func multipartBody(file *multipartFile, raw json.RawMessage) ([]byte, string, string) {
	var in uploadBody

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&in)
	if err != nil {
		return nil, "", "the body argument must be an object with fileName, contentBase64 and optionally contentType"
	}

	if in.FileName == "" {
		return nil, "", "body.fileName is required"
	}

	if in.ContentBase64 == "" {
		return nil, "", "body.contentBase64 is required"
	}

	if len(in.ContentBase64) > maxUploadBase64() {
		return nil, "", fmt.Sprintf("the file is too large: at most %d bytes (%d MiB) are accepted, as on the website", assetstore.MaxUploadBytes, assetstore.MaxUploadBytes/bytesPerMiB)
	}

	content, err := base64.StdEncoding.DecodeString(in.ContentBase64)
	if err != nil {
		return nil, "", "body.contentBase64 is not valid base64 (standard alphabet, padded)"
	}

	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", mime.FormatMediaType("form-data", map[string]string{"name": file.Field, "filename": in.FileName}))

	if in.ContentType != "" {
		header.Set("Content-Type", in.ContentType)
	}

	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	part, err := writer.CreatePart(header)
	if err != nil {
		return nil, "", "the file does not make a request"
	}
	// Writes to a bytes.Buffer cannot fail.
	_, _ = part.Write(content)
	_ = writer.Close()

	return buf.Bytes(), writer.FormDataContentType(), ""
}
