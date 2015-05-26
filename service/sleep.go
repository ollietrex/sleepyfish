//sleep.go
package service

import (
	"github.com/ollietrex/sleepyfish/database"
)

//The service sleep object to deal with the loading and persistence of sleep
type Sleep struct {
	Id        int64
	PersonId  int64
	Start     int64
	End       int64
	Quality   int8
	Feeling   int8
	Comment   string
	CreatedOn int64
}

//Base implementation of the sleep service interface
type sleepService struct {
	sleepRepository database.ISleepRepository
}

//Get a item of sleep by its unique identifier from the persistent storage
func (s sleepService) GetById(id int64, personId int64) (*Sleep, error) {
	dbSleep, err := s.sleepRepository.GetById(id, personId)
	if err != nil {
		return nil, err
	}
	return s.fromDb(dbSleep), nil
}

//Get a list of sleep items from the persistent storage by the unique person identifier
func (s sleepService) GetByPerson(personId int64, start int64, end int64, min int64, max int64) ([]*Sleep, error) {
	dbSleepCollection, err := s.sleepRepository.GetByPerson(personId, start, end, min, max)

	if err != nil {
		return nil, err
	}

	sleepCollection := []*Sleep{}

	for _, dbSleep := range dbSleepCollection {
		sleepCollection = append(sleepCollection, s.fromDb(&dbSleep))
	}

	return sleepCollection, nil
}

//Persist an item of sleep to the persistent storage
func (s sleepService) Persist(sleep *Sleep) (*Sleep, error) {

	dbSleep := &database.Sleep{
		Id:        sleep.Id,
		PersonId:  sleep.PersonId,
		Start:     sleep.Start,
		End:       sleep.End,
		Quality:   sleep.Quality,
		Feeling:   sleep.Feeling,
		Comment:   sleep.Comment,
		CreatedOn: sleep.CreatedOn,
	}

	dbSleep, err := s.sleepRepository.Persist(dbSleep)

	if err != nil {
		return nil, err
	}

	return s.fromDb(dbSleep), nil
}

//Get the count of sleep in the system
func (s sleepService) Count() int64 {
	return s.sleepRepository.Count()
}

//Delete an item item of sleep from the persistent storage
func (s sleepService) Delete(sleep *Sleep) error {
	dbSleep, err := s.sleepRepository.GetById(sleep.Id, sleep.PersonId)
	if err != nil {
		return err
	}
	return s.sleepRepository.Delete(dbSleep)
}

func (s sleepService) fromDb(dbSleep *database.Sleep) *Sleep {
	return &Sleep{
		Id:        dbSleep.Id,
		PersonId:  dbSleep.PersonId,
		Start:     dbSleep.Start,
		End:       dbSleep.End,
		Quality:   dbSleep.Quality,
		Feeling:   dbSleep.Feeling,
		Comment:   dbSleep.Comment,
		CreatedOn: dbSleep.CreatedOn,
	}
}

//The interface definition for sleep
type ISleepService interface {
	//Get a item of sleep by its unique identifier from the persistent storage
	GetById(id int64, personId int64) (*Sleep, error)
	//Get a list of sleep items from the persistent storage by the unique person identifier
	//with the ability to search by date and with paging
	GetByPerson(personId int64, start int64, end int64, min int64, max int64) ([]*Sleep, error)
	//Persist an item of sleep to the persistent storage
	Persist(sleep *Sleep) (*Sleep, error)
	//Delete an item item of sleep from the persistent storage
	Delete(sleep *Sleep) error
	//Provide a count of the amount of sleep recorded in the session
	Count() int64
}

//Get the sleep service for dealing with the persistence of sleep
//into the persistent storage
func GetSleepService() (ISleepService, error) {
	if sleepSrv == nil {
		repo, err := database.GetSleepRepository()
		if err != nil {
			return nil, err
		}
		sleepSrv = &sleepService{
			sleepRepository: repo,
		}
	}
	return sleepSrv, nil
}
