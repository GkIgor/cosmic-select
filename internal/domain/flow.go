package domain

import "errors"

// Phase represents the current step in the core user flow.
type Phase uint8

const (
	PhaseIdle Phase = iota
	PhaseSelecting
	PhaseRecognizing
	PhaseReady
	PhaseRunningAction
)

// Action is one of the four actions supported by the first release.
type Action string

const (
	ActionTranslate Action = "translate"
	ActionCopy      Action = "copy"
	ActionSearch    Action = "search"
	ActionAskAI     Action = "ask_ai"
)

var (
	ErrInvalidTransition = errors.New("invalid application state transition")
	ErrNoText            = errors.New("no text detected")
	ErrInvalidAction     = errors.New("unsupported action")
)

// Session contains only the user-facing state needed after OCR. Image bytes
// are deliberately not retained here; screenshots are transient inputs.
type Session struct {
	Phase  Phase
	Text   string
	Action Action
}

func (s *Session) Activate() error {
	if s.Phase != PhaseIdle {
		return ErrInvalidTransition
	}
	s.Phase = PhaseSelecting
	return nil
}

func (s *Session) BeginRecognition() error {
	if s.Phase != PhaseSelecting {
		return ErrInvalidTransition
	}
	s.Phase = PhaseRecognizing
	return nil
}

func (s *Session) SetText(text string) error {
	if s.Phase != PhaseRecognizing {
		return ErrInvalidTransition
	}
	if text == "" {
		return ErrNoText
	}
	s.Text = text
	s.Phase = PhaseReady
	return nil
}

func (s *Session) ChooseAction(action Action) error {
	if s.Phase != PhaseReady || s.Text == "" {
		return ErrInvalidTransition
	}
	if !IsAction(action) {
		return ErrInvalidAction
	}
	s.Action = action
	s.Phase = PhaseRunningAction
	return nil
}

func (s *Session) Cancel() {
	*s = Session{Phase: PhaseIdle}
}

func IsAction(action Action) bool {
	switch action {
	case ActionTranslate, ActionCopy, ActionSearch, ActionAskAI:
		return true
	default:
		return false
	}
}
