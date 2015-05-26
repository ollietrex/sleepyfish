//session.go
package database

import (
	"gopkg.in/gorp.v1"
	"time"
)

type Session struct {
	Id             int64  `db:"id"`
	Guid           string `db:"guid"`
	PersonId       int64  `db:"person_id"`
	LastActivityOn int64  `db:"last_activity_on"`
	UpdatedOn      int64  `db:"updated_on"`
	CreatedOn      int64  `db:"created_on"`
}

//Pre insert hook to updated on and created on date for an insert
func (i *Session) PreInsert(s gorp.SqlExecutor) error {
	i.UpdatedOn = time.Now().Unix()
	i.CreatedOn = i.UpdatedOn
	return nil
}

//Pre update hook to update the updated date for an update
func (i *Session) PreUpdate(s gorp.SqlExecutor) error {
	i.UpdatedOn = time.Now().Unix()
	return nil
}

//The base implementation of the session repository
type sessionRepository struct {
	db *gorp.DbMap
}

//Get the session by its guid identifier
func (s sessionRepository) GetByGuid(guid string) (*Session, error) {
	dbSession := Session{}
	err := s.db.SelectOne(&dbSession, "SELECT * FROM session WHERE guid=:guid", map[string]interface{}{"guid": guid})
	if err != nil {
		return nil, err
	}
	return &dbSession, nil
}

//Persist the session into the persistent storage
func (s sessionRepository) Persist(session *Session) (*Session, error) {
	var err error
	if session.Id == 0 {
		err = s.db.Insert(session)
	} else {
		_, err = s.db.Update(session)
	}
	if err != nil {
		return nil, err
	}
	return session, nil
}

//Delete a session from the persistent storage
func (s sessionRepository) Delete(session *Session) error {
	_, err := s.db.Delete(session)
	return err
}

//The interface definition for the persistence of sessions
//into the persistent storage
type ISessionRepository interface {
	GetByGuid(guid string) (*Session, error)
	Persist(session *Session) (*Session, error)
	Delete(session *Session) error
}

//Get the session repository to enable the persistence of sessions
//into the persistent storage
func GetSessionRepository() (ISessionRepository, error) {
	if sessionRepo == nil {
		db, err := GetDatabase()
		if err != nil {
			return nil, err
		}
		sessionRepo = &sessionRepository{
			db: db,
		}
	}
	return sessionRepo, nil
}
