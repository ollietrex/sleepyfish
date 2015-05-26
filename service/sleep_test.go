//sleep_test.go
package service

import (
	"testing"
	"time"
)

func TestSleepService(t *testing.T) {
	person, _, err := getTestPerson()

	if err != nil {
		t.Error(err)
	}

	sleepService, err := GetSleepService()

	if err != nil {
		t.Error(err)
	}

	newSleep := &Sleep{
		PersonId: person.Id,
		Start:    time.Now().Unix(),
		End:      time.Now().Unix(),
		Quality:  10,
		Feeling:  10,
		Comment:  "This is a test sleep",
	}

	newSleep, err = sleepService.Persist(newSleep)

	if err != nil {
		t.Error(err)
	}

	if newSleep.Id == 0 {
		t.Error("The sleep did not persist to the database")
	}

	testSleep, err := sleepService.GetById(newSleep.Id, person.Id)

	if newSleep.Comment != testSleep.Comment {
		t.Error("Comment does not match")
	}

	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	sleepCollection, err := sleepService.GetByPerson(person.Id, start.Unix(), end.Unix(), 0, 10)

	if err != nil {
		t.Error(err)
	}

	found := false

	t.Log(sleepCollection)

	for _, sleep := range sleepCollection {
		if sleep.Id == testSleep.Id {
			found = true
		}
	}

	if !found {
		t.Error("Sleep not found in collection")
	}

	if err != nil {
		t.Error(err)
	}

	err = sleepService.Delete(testSleep)

	if err != nil {
		t.Error(err)
	}

}

func persistTestSleep(person *Person) (*Sleep, error) {
	sleepService, err := GetSleepService()

	if err != nil {
		return nil, err
	}

	newSleep := &Sleep{
		PersonId: person.Id,
		Start:    time.Now().Unix(),
		End:      time.Now().Unix(),
		Quality:  10,
		Feeling:  10,
		Comment:  "This is a test sleep",
	}

	newSleep, err = sleepService.Persist(newSleep)

	if err != nil {
		return nil, err
	}

	return newSleep, nil
}

func TestSleepCount(t *testing.T) {

	person, _, err := getTestPerson()

	if err != nil {
		t.Error(err)
	}

	for i := 1; i <= 10; i++ {
		persistTestSleep(person)
	}

	sleepService, err := GetSleepService()

	if err != nil {
		t.Error(err)
	}

	count := sleepService.Count()

	if count < 10 {
		t.Error("Not enough sleep found in the system")
	}

}
