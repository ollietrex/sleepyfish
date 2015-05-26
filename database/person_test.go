//person_test.go
package database

import (
	"github.com/ollietrex/sleepyfish/helpers"
	"testing"
)

func getTestPerson() *Person {
	testPerson := &Person{
		Name:         helpers.RandomName(),
		EmailAddress: helpers.RandomEmail(),
		PasswordHash: []byte(helpers.RandomString(32)),
		PasswordSalt: []byte(helpers.RandomString(32)),
	}
	return testPerson
}

func persistTestPerson() (*Person, error) {
	testPerson := getTestPerson()

	personRepo, err := GetPersonRepository()

	if err != nil {
		return nil, err
	}

	testPerson, err = personRepo.Persist(testPerson)

	if err != nil {
		return nil, err
	}

	return testPerson, nil
}

func TestPersonRepository(t *testing.T) {
	personRepo, err := GetPersonRepository()

	if err != nil {
		t.Error(err)
	}
	//Test if we can insert some one into the database

	testPerson := getTestPerson()

	testPerson, err = personRepo.Persist(testPerson)

	if err != nil {
		t.Error(err)
	}

	if testPerson.Id == 0 {
		t.Error("Id has not been set")
	}

	//Test if we can get them back out by id

	byIdPerson, err := personRepo.GetById(testPerson.Id)

	if err != nil {
		t.Error(err)
	}

	if byIdPerson.Name != testPerson.Name {
		t.Error("The name of the two people is not the same")
	}

	//Test if we can get a person by email address

	emailPerson, err := personRepo.GetByEmail(testPerson.EmailAddress)

	if err != nil {
		t.Error(err)
	}

	if emailPerson.Name != testPerson.Name {
		t.Error("The name of the two people is not the same")
	}

	//Delete the person so we dont clog up the db

	err = personRepo.Delete(testPerson)

	if err != nil {
		t.Error(err)
	}

}

func TestPersonRepositoryCount(t *testing.T) {

	for i := 1; i <= 10; i++ {
		persistTestPerson()
	}

	personRepo, err := GetPersonRepository()

	if err != nil {
		t.Error(err)
	}
	count := personRepo.Count()

	if count < 10 {
		t.Error("Not enough people in the system")
	}
}
