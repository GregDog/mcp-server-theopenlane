package tools

import (
	"context"
	"testing"
)

func TestDeleteRequiresID(t *testing.T) {
	h := &handlers{api: &fakeAPI{}}

	_, _, err := h.deleteTask(context.Background(), nil, getInput{})
	if err != errIDRequired {
		t.Fatalf("delete task: got %v", err)
	}
}
