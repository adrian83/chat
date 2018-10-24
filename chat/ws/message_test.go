package ws

import "testing"

func TestUnknownMessageImplementsDoer(t *testing.T) {
	var _ Doer = &UnknownMessage{}
}
