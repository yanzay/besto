package tbot

import "testing"

func TestAccessDenied(t *testing.T) {
	message := mockMessage()
	go accessDenied(message)
	reply := <-message.Replies
	if reply == nil {
		t.Fail()
	}
}

func TestAuthSuccess(t *testing.T) {
	auth := NewAuth([]string{"me"})
	message := mockMessage()
	invoked := false
	handler := func(m *Message) { invoked = true }
	auth(handler)(message)
	if !invoked {
		t.Fail()
	}
}

func TestAuthFail(t *testing.T) {
	auth := NewAuth([]string{"notme"})
	message := mockMessage()
	invoked := false
	handler := func(m *Message) { invoked = true }
	go auth(handler)(message)
	<-message.Replies
	if invoked {
		t.Fail()
	}
}
