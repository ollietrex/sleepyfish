//sleep_test.go
package database

import (
	"testing"
	"time"
)

func TestSleepRepository(t *testing.T) {

	person, err := persistTestPerson()

	if err != nil {
		t.Error(err)
	}

	sleepRepo, err := GetSleepRepository()

	if err != nil {
		t.Error(err)
	}

	testSleep := &Sleep{
		PersonId: person.Id,
		Start:    time.Now().Unix(),
		End:      time.Now().Unix(),
		Quality:  10,
		Feeling:  10,
		Comment:  "Hello",
	}

	testSleep, err = sleepRepo.Persist(testSleep)

	if err != nil {
		t.Error(err)
	}

	if testSleep.Id == 0 {
		t.Error("Id has not been set")
	}

	newSleep, err := sleepRepo.GetById(testSleep.Id, person.Id)

	if err != nil {
		t.Error(err)
	}

	if testSleep.Start != newSleep.Start {
		t.Error("The starts times are not the same for both records")
	}

	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	sleepCollection, err := sleepRepo.GetByPerson(person.Id, start.Unix(), end.Unix(), 0, 10)

	found := false
	for _, item := range sleepCollection {
		if item.Id == testSleep.Id {
			found = true
		}
	}

	if !found {
		t.Error("The item was not found in the collection")
	}

	err = sleepRepo.Delete(testSleep)

	if err != nil {
		t.Error(err)
	}

}

func persistTestSleep(person *Person) (*Sleep, error) {
	sleepRepo, err := GetSleepRepository()

	if err != nil {
		return nil, err
	}

	testSleep := &Sleep{
		PersonId: person.Id,
		Start:    time.Now().Unix(),
		End:      time.Now().Unix(),
		Quality:  10,
		Feeling:  10,
		Comment:  "Hello",
	}

	testSleep, err = sleepRepo.Persist(testSleep)

	if err != nil {
		return nil, err
	}

	return testSleep, nil
}

func TestSleepRepositoryCount(t *testing.T) {

	person, err := persistTestPerson()

	if err != nil {
		t.Error(err)
	}

	for i := 1; i <= 10; i++ {
		persistTestSleep(person)
	}

	sleepRepo, err := GetSleepRepository()

	if err != nil {
		t.Error(err)
	}

	count := sleepRepo.Count()

	if count < 10 {
		t.Error("Not enough sleep in the system")
	}
}
