package tools

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func TestCreateEntityRequiresName(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	_, _, err := h.createEntity(context.Background(), nil, createEntityInput{})
	if err != errNameRequired {
		t.Fatalf("create entity: got %v", err)
	}
}

func TestUpdateEntityRequiresFields(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	_, _, err := h.updateEntity(context.Background(), nil, updateEntityInput{ID: "ent_1"})
	if err != errUpdateFieldsRequired {
		t.Fatalf("update entity: got %v", err)
	}
}

func TestUpdateEntityAllowsLogoOnly(t *testing.T) {
	h := &handlers{
		api:            &fakeAPI{entity: &openlane.EntityDetail{ID: "ent_1", LogoFileID: strPtr("file_1")}},
		allowWrite:     true,
		maxUploadBytes: openlane.DefaultMaxUploadBytes,
	}
	_, item, err := h.updateEntity(context.Background(), nil, updateEntityInput{
		ID: "ent_1",
		Logo: &entityLogoInput{
			Filename:      "logo.png",
			ContentType:   "image/png",
			ContentBase64: base64.StdEncoding.EncodeToString([]byte("png")),
		},
	})
	if err != nil {
		t.Fatalf("update entity with logo: %v", err)
	}
	if item.LogoFileID != "file_1" {
		t.Fatalf("logo file id: got %q", item.LogoFileID)
	}
}

func TestDecodeEntityLogoRejectsMissingFilename(t *testing.T) {
	h := &handlers{maxUploadBytes: openlane.DefaultMaxUploadBytes}
	_, err := h.decodeEntityLogo(&entityLogoInput{ContentBase64: "dGVzdA=="})
	if err == nil || !strings.Contains(err.Error(), "filename is required") {
		t.Fatalf("expected filename error, got %v", err)
	}
}

func TestDecodeEntityLogoRejectsOversize(t *testing.T) {
	h := &handlers{maxUploadBytes: 4}
	_, err := h.decodeEntityLogo(&entityLogoInput{
		Filename:      "logo.png",
		ContentBase64: base64.StdEncoding.EncodeToString([]byte("12345")),
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds maximum upload size") {
		t.Fatalf("expected oversize error, got %v", err)
	}
}

func TestCreateEntitySuccess(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	_, item, err := h.createEntity(context.Background(), nil, createEntityInput{Name: "Acme Corp"})
	if err != nil {
		t.Fatalf("create entity: %v", err)
	}
	if item.ID != "ent_1" || item.Name != "Acme Corp" {
		t.Fatalf("unexpected item: %+v", item)
	}
}
