package openlane

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestDecodeEvidenceUploads(t *testing.T) {
	content := base64.StdEncoding.EncodeToString([]byte("hello"))
	uploads, err := DecodeEvidenceUploads([]EvidenceFile{{
		Filename:      "note.txt",
		ContentType:   "text/plain",
		ContentBase64: content,
	}}, DefaultMaxUploadBytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(uploads) != 1 || uploads[0].Filename != "note.txt" || uploads[0].Size != 5 {
		t.Fatalf("unexpected uploads: %+v", uploads)
	}
}

func TestDecodeEvidenceUploadsRejectsOversize(t *testing.T) {
	large := base64.StdEncoding.EncodeToString([]byte(strings.Repeat("a", 32)))
	_, err := DecodeEvidenceUploads([]EvidenceFile{{
		Filename:      "big.bin",
		ContentBase64: large,
	}}, 16)
	if err == nil {
		t.Fatal("expected size error")
	}
}
