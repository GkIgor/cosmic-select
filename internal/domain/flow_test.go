package domain

import "testing"

func TestSessionFlow(t *testing.T) {
	var session Session

	steps := []struct {
		name string
		step func() error
	}{
		{name: "activate", step: session.Activate},
		{name: "recognize", step: session.BeginRecognition},
		{name: "set text", step: func() error { return session.SetText("visible text") }},
		{name: "choose action", step: func() error { return session.ChooseAction(ActionCopy) }},
	}
	for _, step := range steps {
		if err := step.step(); err != nil {
			t.Fatalf("%s: %v", step.name, err)
		}
	}

	if session.Phase != PhaseRunningAction || session.Text != "visible text" {
		t.Fatalf("unexpected session state: %+v", session)
	}
}

func TestSessionRejectsEmptyOCRResult(t *testing.T) {
	session := Session{Phase: PhaseRecognizing}
	if err := session.SetText(""); err != ErrNoText {
		t.Fatalf("expected ErrNoText, got %v", err)
	}
}

func TestCancelClearsTransientState(t *testing.T) {
	session := Session{Phase: PhaseReady, Text: "private text", Action: ActionCopy}
	session.Cancel()

	if session != (Session{Phase: PhaseIdle}) {
		t.Fatalf("session was not cleared: %+v", session)
	}
}
