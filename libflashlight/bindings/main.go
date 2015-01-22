// package flashlight provides minimal configuration for spawning a flashlight
// client.

package flashlight

import (
	"github.com/getlantern/lantern-android/client"
	"strings"
)

var DefaultClient *client.Client

// StopClientProxy stops the proxy.
func StopClientProxy() error {
	DefaultClient.Stop()
	return nil
}

// RunClientProxy creates a new client at the given address.
func RunClientProxy(listenAddr string) error {

	DefaultClient = client.NewClient(listenAddr)

	go func() {
		var err error
		if err = DefaultClient.ListenAndServe(); err != nil {
			// Error is not exported: https://golang.org/src/net/net.go#L284
			if !strings.Contains(err.Error(), "use of closed network connection") {
				panic(err.Error())
			}
		}
	}()
	return nil
}
