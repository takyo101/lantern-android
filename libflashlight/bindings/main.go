// package flashlight provides minimal configuration for spawning a flashlight
// client.

package flashlight

import (
	"github.com/getlantern/lantern-android/client"
)

var DefaultClient *client.Client

// StopClientProxy stops the proxy.
func StopClientProxy() error {
	return DefaultClient.Stop()
}

// RunClientProxy creates a new client at the given address. If an active
// client is found it kill the client before starting a new one.
func RunClientProxy(listenAddr string) (err error) {
	DefaultClient = client.NewClient(listenAddr)
	return nil
}
