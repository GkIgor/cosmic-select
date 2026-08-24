package app

import (
	"context"

	"github.com/GkIgor/cosmic-select/internal/domain"
	"github.com/GkIgor/cosmic-select/internal/ocr"
	"github.com/GkIgor/cosmic-select/internal/ports"
)

// Coordinator owns the product flow without depending on GTK or D-Bus.
type Coordinator struct {
	selector ports.Selector
	ocr      ports.OCREngine
	session  domain.Session
}

func NewCoordinator(selector ports.Selector, engine ports.OCREngine) *Coordinator {
	return &Coordinator{selector: selector, ocr: engine}
}

func (c *Coordinator) Session() domain.Session {
	return c.session
}

// HandleActivation runs the transient part of the core flow. The selected
// image is passed directly from capture to OCR and is never stored by the
// coordinator.
func (c *Coordinator) HandleActivation(ctx context.Context) error {
	if err := c.session.Activate(); err != nil {
		return err
	}

	image, err := c.selector.Select(ctx)
	if err != nil {
		c.session.Cancel()
		return err
	}
	if err := c.session.BeginRecognition(); err != nil {
		c.session.Cancel()
		return err
	}

	text, err := c.ocr.Extract(ctx, image)
	if err != nil {
		c.session.Cancel()
		return err
	}
	if err := c.session.SetText(ocr.CleanText(text)); err != nil {
		c.session.Cancel()
		return err
	}

	return nil
}

func (c *Coordinator) Cancel() {
	c.session.Cancel()
}
