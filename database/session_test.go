//session_test.go
package database

import (
	"testing"
	"time"
)

func TestSessionRepository(t *testing.T) {

	sessionRepo, err := GetSessionRepository()

	if err != nil {
		t.Error(err)
	}

	guid := "aaaaaa-aaaaaa-aaaaaa-aaaaaa-aaaaaa"

	testSession := &Session{
		Guid:           guid,
		PersonId:       1,
		LastActivityOn: time.Now().Unix(),
	}

	testSession, err = sessionRepo.Persist(testSession)

	if err != nil {
		t.Error(err)
	}

	if testSession.Id == 0 {
		t.Error("Did not create the session record")
	}

	byGuidSession, err := sessionRepo.GetByGuid(guid)

	if err != nil {
		t.Error(err)
	}

	if testSession.Id != byGuidSession.Id {
		t.Error("The same session was not loaded")
	}

	lastActivityOn := time.Now().Unix() + 100

	testSession.LastActivityOn = lastActivityOn

	testSession, err = sessionRepo.Persist(testSession)

	if err != nil {
		t.Error(err)
	}

	byGuidSession, err = sessionRepo.GetByGuid(guid)

	if err != nil {
		t.Error(err)
	}

	if byGuidSession.LastActivityOn != lastActivityOn {
		t.Error("Session did not update")
	}

	err = sessionRepo.Delete(testSession)

	if err != nil {
		t.Error(err)
	}

}
