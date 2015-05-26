//session_test.go
package service

import (
	"testing"
)

func TestSessionService(t *testing.T) {

	person, _, err := getTestPerson()

	if err != nil {
		t.Error(err)
	}

	sessionService, err := GetSessionService()

	session, err := sessionService.RegisterSession(person.Id)

	if err != nil {
		t.Error(err)
	}

	if session.Id == 0 {
		t.Error("Session has not been created")
	}

	newSession, err := sessionService.GetSession(session.Guid)

	if err != nil {
		t.Error(err)
	}

	if newSession.Guid != session.Guid {
		t.Error("Session was not created correctly")
	}

	err = sessionService.DeleteSession(session.Guid)

	_, err = sessionService.GetSession(session.Guid)

	if err == nil {
		t.Error("The session was not deleted")
	}

}
