package client

import (
	"context"
	"testing"
	"time"
)

func TestController_NextTurn(t *testing.T) {
	ctx, stop := context.WithCancel(context.Background())

	defer stop()

	_, ctrl, err := NewTestController(ctx, "localhost", "8080", "local")
	if err != nil {
		t.Fatal(err.Error())
	}

	ctrl.NextTurn()
	time.Sleep(5 * time.Second)
}
