package portal

import (
	"context"
	"fmt"
	"sync"

	"github.com/godbus/dbus/v5"
)

const (
	desktopBusName = "org.freedesktop.portal.Desktop"
	desktopPath    = dbus.ObjectPath("/org/freedesktop/portal/desktop")
)

// Client owns the session-bus connection used by the portal adapters.
type Client struct {
	conn *dbus.Conn
	mu   sync.Mutex
}

func NewClient() (*Client, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("connect to session bus: %w", err)
	}
	return &Client{conn: conn}, nil
}

func newClient(conn *dbus.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) object() dbus.BusObject {
	return c.conn.Object(desktopBusName, desktopPath)
}

func (c *Client) waitResponse(ctx context.Context, request dbus.ObjectPath) (map[string]dbus.Variant, error) {
	if err := c.conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.portal.Request"),
		dbus.WithMatchMember("Response"),
		dbus.WithMatchObjectPath(request),
	); err != nil {
		return nil, fmt.Errorf("subscribe to portal response: %w", err)
	}
	defer c.conn.RemoveMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.portal.Request"),
		dbus.WithMatchMember("Response"),
		dbus.WithMatchObjectPath(request),
	)

	signals := make(chan *dbus.Signal, 1)
	c.conn.Signal(signals)
	defer c.conn.RemoveSignal(signals)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case signal := <-signals:
		if len(signal.Body) != 2 {
			return nil, ErrPortalResponse
		}
		response, ok := signal.Body[0].(uint32)
		if !ok {
			return nil, ErrPortalResponse
		}
		if response == 1 {
			return nil, ErrPortalCancelled
		}
		if response != 0 {
			return nil, fmt.Errorf("%w: response code %d", ErrPortalResponse, response)
		}
		results, ok := signal.Body[1].(map[string]dbus.Variant)
		if !ok {
			return nil, ErrPortalResponse
		}
		return results, nil
	}
}

func variantString(results map[string]dbus.Variant, key string) (string, error) {
	value, ok := results[key]
	if !ok {
		return "", fmt.Errorf("portal response missing %q", key)
	}
	result, ok := value.Value().(string)
	if !ok {
		return "", fmt.Errorf("portal response %q has unexpected type", key)
	}
	return result, nil
}
