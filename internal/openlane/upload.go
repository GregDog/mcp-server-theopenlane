package openlane

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

const DefaultMaxUploadBytes = 10 * 1024 * 1024 // 10 MiB decoded

// EvidenceFile is a file attached to evidence create/update.
type EvidenceFile struct {
	Filename      string
	ContentType   string
	ContentBase64 string
}

// DecodeEvidenceUploads decodes base64 file payloads into GraphQL uploads.
func DecodeEvidenceUploads(files []EvidenceFile, maxBytes int64) ([]*graphql.Upload, error) {
	if len(files) == 0 {
		return nil, nil
	}
	if maxBytes <= 0 {
		maxBytes = DefaultMaxUploadBytes
	}

	uploads := make([]*graphql.Upload, 0, len(files))
	for i, f := range files {
		if strings.TrimSpace(f.Filename) == "" {
			return nil, fmt.Errorf("files[%d]: filename is required", i)
		}
		if strings.TrimSpace(f.ContentBase64) == "" {
			return nil, fmt.Errorf("files[%d]: content_base64 is required", i)
		}

		raw, err := base64.StdEncoding.DecodeString(f.ContentBase64)
		if err != nil {
			return nil, fmt.Errorf("files[%d]: invalid base64: %w", i, err)
		}
		if int64(len(raw)) > maxBytes {
			return nil, fmt.Errorf("files[%d]: exceeds maximum upload size of %d bytes", i, maxBytes)
		}

		uploads = append(uploads, &graphql.Upload{
			File:        bytes.NewReader(raw),
			Filename:    f.Filename,
			Size:        int64(len(raw)),
			ContentType: strings.TrimSpace(f.ContentType),
		})
	}
	return uploads, nil
}

// UploadReader returns the upload reader, rewinding when possible.
func UploadReader(u *graphql.Upload) (io.Reader, error) {
	if u == nil || u.File == nil {
		return nil, fmt.Errorf("upload file is required")
	}
	if seeker, ok := u.File.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("rewind upload: %w", err)
		}
	}
	return u.File, nil
}
