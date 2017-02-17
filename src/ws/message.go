package ws

import (
	"errors"
)

type Message map[string]interface{}

func (m Message) getType() (string, error) {
	iType, ok := m["msgType"]
	if !ok {
		return "", errors.New("No property 'msgType'")
	}
	sType, ok := iType.(string)
	if !ok {
		return "", errors.New("Property 'msgType'is not a string")
	}
	return sType, nil
}
