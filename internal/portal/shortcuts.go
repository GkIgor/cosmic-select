package portal

import (
	"context"
	"fmt"
	"sync"

	"github.com/godbus/dbus/v5"
)

const globalShortcutsInterface = "org.freedesktop.portal.GlobalShortcuts"

type shortcutRegistration struct {
	ID      string
	Options map[string]dbus.Variant
}

// GlobalShortcuts registers one application activation shortcut through the
// XDG Global Shortcuts portal.
type GlobalShortcuts struct {
	client     *Client
	session    dbus.ObjectPath
	shortcutID string
	stop       chan struct{}
	stopped    chan struct{}
	signals    chan *dbus.Signal
	closeOnce  sync.Once
	registered bool
	callback   func()
}

func NewGlobalShortcuts(client *Client, shortcutID string) *GlobalShortcuts {
	return &GlobalShortcuts{client: client, shortcutID: shortcutID}
}

func (s *GlobalShortcuts) Register(ctx context.Context, trigger func()) error {
	s.client.mu.Lock()
	defer s.client.mu.Unlock()
	if s.registered {
		return fmt.Errorf("global shortcut is already registered")
	}

	createCall := s.client.object().CallWithContext(ctx, globalShortcutsInterface+".CreateSession", 0,
		map[string]dbus.Variant{
			"handle_token":         dbus.MakeVariant("cosmic_select_session"),
			"session_handle_token": dbus.MakeVariant("cosmic_select"),
		},
	)
	if createCall.Err != nil {
		return fmt.Errorf("create global shortcut session: %w", createCall.Err)
	}
	var createRequest dbus.ObjectPath
	if err := createCall.Store(&createRequest); err != nil {
		return fmt.Errorf("read global shortcut request handle: %w", err)
	}
	created, err := s.client.waitResponse(ctx, createRequest)
	if err != nil {
		return err
	}
	session, err := variantObjectPath(created, "session_handle")
	if err != nil {
		return err
	}

	bindCall := s.client.object().CallWithContext(ctx, globalShortcutsInterface+".BindShortcuts", 0,
		session,
		[]shortcutRegistration{{
			ID: s.shortcutID,
			Options: map[string]dbus.Variant{
				"description":       dbus.MakeVariant("Select text from the screen"),
				"preferred_trigger": dbus.MakeVariant("<Super><Shift>s"),
			},
		}},
		"",
		map[string]dbus.Variant{"handle_token": dbus.MakeVariant("cosmic_select_bind")},
	)
	if bindCall.Err != nil {
		return fmt.Errorf("bind global shortcut: %w", bindCall.Err)
	}
	var bindRequest dbus.ObjectPath
	if err := bindCall.Store(&bindRequest); err != nil {
		return fmt.Errorf("read bind request handle: %w", err)
	}
	if _, err := s.client.waitResponse(ctx, bindRequest); err != nil {
		return err
	}

	s.session = session
	s.callback = trigger
	s.stop = make(chan struct{})
	s.stopped = make(chan struct{})
	s.signals = make(chan *dbus.Signal, 8)
	s.client.conn.Signal(s.signals)
	if err := s.client.conn.AddMatchSignal(
		dbus.WithMatchInterface(globalShortcutsInterface),
		dbus.WithMatchMember("Activated"),
		dbus.WithMatchObjectPath(session),
	); err != nil {
		s.client.conn.RemoveSignal(s.signals)
		return fmt.Errorf("subscribe to global shortcut activation: %w", err)
	}
	s.registered = true
	go s.listen()
	return nil
}

func (s *GlobalShortcuts) listen() {
	defer close(s.stopped)
	for {
		select {
		case <-s.stop:
			return
		case signal := <-s.signals:
			if len(signal.Body) < 2 || signal.Body[0] != s.session || signal.Body[1] != s.shortcutID {
				continue
			}
			if s.callback != nil {
				s.callback()
			}
		}
	}
}

func (s *GlobalShortcuts) Close() error {
	var closeErr error
	s.closeOnce.Do(func() {
		if !s.registered {
			return
		}
		close(s.stop)
		<-s.stopped
		s.client.conn.RemoveSignal(s.signals)
		_ = s.client.conn.RemoveMatchSignal(
			dbus.WithMatchInterface(globalShortcutsInterface),
			dbus.WithMatchMember("Activated"),
			dbus.WithMatchObjectPath(s.session),
		)
		call := s.client.conn.Object(desktopBusName, s.session).Call("org.freedesktop.portal.Session.Close", 0)
		closeErr = call.Err
	})
	return closeErr
}

func variantObjectPath(results map[string]dbus.Variant, key string) (dbus.ObjectPath, error) {
	value, ok := results[key]
	if !ok {
		return "", fmt.Errorf("portal response missing %q", key)
	}
	switch result := value.Value().(type) {
	case dbus.ObjectPath:
		return result, nil
	case string:
		return dbus.ObjectPath(result), nil
	default:
		return "", fmt.Errorf("portal response %q has unexpected type", key)
	}
}
