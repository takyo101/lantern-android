package client

import (
	"testing"
)

func TestListenAndServe(t *testing.T) {
	var err error

	c := NewClient(":9997")

	if err = c.ListenAndServe(); err != nil {
		t.Fatal(err)
	}
}
