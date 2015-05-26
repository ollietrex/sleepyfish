//sleep.go
package database

import (
	"gopkg.in/gorp.v1"
	"time"
)

//The database persistence object for sleep
type Sleep struct {
	Id        int64  `db:"id"`
	PersonId  int64  `db:"person_id"`
	Start     int64  `db:"start"`
	End       int64  `db:"end"`
	Quality   int8   `db:"quality"`
	Feeling   int8   `db:"feeling"`
	Comment   string `db:"comment"`
	UpdatedOn int64  `db:"updated_on"`
	CreatedOn int64  `db:"created_on"`
}

//Pre insert hook to updated on and created on date for an insert
func (i *Sleep) PreInsert(s gorp.SqlExecutor) error {
	i.UpdatedOn = time.Now().Unix()
	i.CreatedOn = i.UpdatedOn
	return nil
}

//Pre update hook to update the updated date for an update
func (i *Sleep) PreUpdate(s gorp.SqlExecutor) error {
	i.UpdatedOn = time.Now().Unix()
	return nil
}

//The base implementation of the sleep repository
type sleepRepository struct {
	db *gorp.DbMap
}

//Get a item of sleep by its unique identifier from the persistent storage
func (s sleepRepository) GetById(id int64, personId int64) (*Sleep, error) {
	sleep := &Sleep{}
	err := s.db.SelectOne(sleep, "SELECT * FROM sleep WHERE id=:id AND person_id=:person_id", map[string]interface{}{"id": id, "person_id": personId})
	if err != nil {
		return nil, err
	}
	return sleep, nil
}

//Get a list of sleep items from the persistent storage by the unique person identifier
func (s sleepRepository) GetByPerson(personId int64, start int64, end int64, min int64, max int64) ([]Sleep, error) {
	var sleep []Sleep

	count := max - min

	filters := map[string]interface{}{"person_id": personId, "start": start, "end": end, "min": min, "count": count}

	_, err := s.db.Select(&sleep, "SELECT * FROM sleep WHERE person_id=:person_id AND start >= :start AND end <= :end ORDER BY start DESC LIMIT :min,:count;", filters)
	if err != nil {
		return nil, err
	}
	return sleep, nil
}

//Persist an item of sleep to the persistent storage
func (s sleepRepository) Persist(sleep *Sleep) (*Sleep, error) {
	var err error
	if sleep.Id == 0 {
		err = s.db.Insert(sleep)
	} else {
		_, err = s.db.Update(sleep)
	}
	if err != nil {
		return nil, err
	}
	return sleep, nil
}

//Delete an item item of sleep from the persistent storage
func (s sleepRepository) Delete(sleep *Sleep) error {
	_, err := s.db.Delete(sleep)
	return err
}

//Count the amount of sleep in the persistent storage
func (s sleepRepository) Count() int64 {
	count, _ := s.db.SelectInt("SELECT COUNT(*) FROM sleep")
	return count
}

//The interface for the persistence of sleep to the database
type ISleepRepository interface {
	//Get a item of sleep by its unique identifier from the persistent storage
	GetById(id int64, personId int64) (*Sleep, error)
	//Get a list of sleep items from the persistent storage by the unique person identifier
	//with a date range with paging
	GetByPerson(personId int64, start int64, end int64, min int64, max int64) ([]Sleep, error)
	//Persist an item of sleep to the persistent storage
	Persist(sleep *Sleep) (*Sleep, error)
	//Delete an item item of sleep from the persistent storage
	Delete(sleep *Sleep) error
	//Count the amount of sleep in the persistent storage
	Count() int64
}

//Get an instance of the sleep repository to deal with the persistence of sleep
func GetSleepRepository() (ISleepRepository, error) {
	if sleepRepo == nil {
		db, err := GetDatabase()
		if err != nil {
			return nil, err
		}
		sleepRepo = &sleepRepository{
			db: db,
		}
	}
	return sleepRepo, nil
}
