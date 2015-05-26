//session.go
package service

import (
	"bitbucket.org/ollietrex/sleepyfish/database"
	"github.com/twinj/uuid"
	"time"
)

//The service session object for managing sessions in the persistent storage
type Session struct {
	Id             int64
	Guid           string
	PersonId       int64
	LastActivityOn int64
}

//Base implementation of the session service interface
type sessionService struct {
	sessionRepository database.ISessionRepository
}

//Register a new session once the user has authentciated in the system
func (s sessionService) RegisterSession(personId int64) (*Session, error) {
	dbSession := &database.Session{
		Guid:           uuid.NewV4().String(),
		PersonId:       personId,
		LastActivityOn: time.Now().Unix(),
	}

	dbSession, err := s.sessionRepository.Persist(dbSession)

	if err != nil {
		return nil, err
	}

	newSession := &Session{
		Id:             dbSession.Id,
		Guid:           dbSession.Guid,
		PersonId:       dbSession.PersonId,
		LastActivityOn: dbSession.LastActivityOn,
	}

	return newSession, nil
}

//Get a session from the persistent storage based on its guid
func (s sessionService) GetSession(guid string) (*Session, error) {
	dbSession, err := s.sessionRepository.GetByGuid(guid)

	if err != nil {
		return nil, err
	}

	newSession := &Session{
		Id:             dbSession.Id,
		Guid:           dbSession.Guid,
		PersonId:       dbSession.PersonId,
		LastActivityOn: dbSession.LastActivityOn,
	}

	return newSession, nil
}

//Refresh the session activity, this should be called on every session usage
//and the service will decide if we are going to refresh the session
func (s sessionService) RefreshSession(session *Session) {
	lastActivityOn := time.Unix(session.LastActivityOn, 0)
	hourAgo := time.Now().Add(-time.Hour)
	if hourAgo.After(lastActivityOn) {
		dbSession, err := s.sessionRepository.GetByGuid(session.Guid)
		if err == nil {
			dbSession.LastActivityOn = time.Now().Unix()
			s.sessionRepository.Persist(dbSession)
		}
	}
}

//Delete a session from the persistent storage
func (s sessionService) DeleteSession(guid string) error {
	dbSession, err := s.sessionRepository.GetByGuid(guid)
	if err != nil {
		return err
	}
	return s.sessionRepository.Delete(dbSession)
}

//The session service interface to deal with the session logic
type ISessionService interface {
	//Register a new session once the user has authentciated in the system
	RegisterSession(personId int64) (*Session, error)
	//Get a session from the persistent storage based on its guid
	GetSession(guid string) (*Session, error)
	//Refresh the session activity, this should be called on every session usage
	//and the service will decide if we are going to refresh the session
	RefreshSession(session *Session)
	//Delete a session from the persistent storage
	DeleteSession(guid string) error
}

func GetSessionService() (ISessionService, error) {
	if sessionSrv == nil {
		repo, err := database.GetSessionRepository()
		if err != nil {
			return nil, err
		}
		sessionSrv = &sessionService{
			sessionRepository: repo,
		}
	}
	return sessionSrv, nil
}
