package app

import (
	"context"
	"errors"
	"testing"

	"github.com/GkIgor/cosmic-select/internal/domain"
	"github.com/GkIgor/cosmic-select/internal/ports"
)

type selectorStub struct {
	image ports.Image
	err   error
}

func (s selectorStub) Select(context.Context) (ports.Image, error) {
	return s.image, s.err
}

type ocrStub struct {
	text string
	err  error
}

func (s ocrStub) Extract(context.Context, ports.Image) (string, error) {
	return s.text, s.err
}

func TestHandleActivationProducesReadySession(t *testing.T) {
	coordinator := NewCoordinator(
		selectorStub{image: ports.Image{Data: []byte("temporary")}},
		ocrStub{text: "  Texto reconhecido  \n"},
	)

	if err := coordinator.HandleActivation(context.Background()); err != nil {
		t.Fatalf("HandleActivation() error = %v", err)
	}

	session := coordinator.Session()
	if session.Phase != domain.PhaseReady || session.Text != "Texto reconhecido" {
		t.Fatalf("unexpected session: %+v", session)
	}
}

func TestHandleActivationCancelsAfterCaptureFailure(t *testing.T) {
	coordinator := NewCoordinator(selectorStub{err: errors.New("portal unavailable")}, ocrStub{})
	if err := coordinator.HandleActivation(context.Background()); err == nil {
		t.Fatal("expected capture error")
	}

	if session := coordinator.Session(); session.Phase != domain.PhaseIdle {
		t.Fatalf("expected idle session after failure, got %+v", session)
	}
}

func TestHandleActivationCancelsWhenOCRFindsNoText(t *testing.T) {
	coordinator := NewCoordinator(selectorStub{}, ocrStub{text: " \n "})
	if err := coordinator.HandleActivation(context.Background()); err != domain.ErrNoText {
		t.Fatalf("expected ErrNoText, got %v", err)
	}

	if session := coordinator.Session(); session.Phase != domain.PhaseIdle {
		t.Fatalf("expected idle session after empty OCR, got %+v", session)
	}
}
