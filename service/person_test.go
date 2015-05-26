//person_test.go
package service

import (
	"bitbucket.org/ollietrex/sleepyfish/helpers"
	"testing"
)

func getTestPerson() (*Person, string, error) {
	service, err := GetPersonService()

	if err != nil {
		return nil, "", err
	}

	password := helpers.RandomString(10)

	signupPerson := &Person{
		Name:         helpers.RandomName(),
		EmailAddress: helpers.RandomEmail(),
	}

	newPerson, err := service.Register(signupPerson, password)

	if err != nil {
		return nil, "", err
	}

	return newPerson, password, nil
}

func TestSignupAndLogin(t *testing.T) {
	service, err := GetPersonService()

	if err != nil {
		t.Error("Error creating the person service", err)
	}

	password := helpers.RandomString(10)

	signupPerson := &Person{
		Name:         helpers.RandomName(),
		EmailAddress: helpers.RandomEmail(),
	}

	newPerson, err := service.Register(signupPerson, password)

	if err != nil {
		t.Error("Error registering", err)
	}

	_, err = service.Register(signupPerson, password)

	if err == nil {
		t.Error("We should have a report of a duplicate account")
	}

	checkPerson(t, newPerson, signupPerson)

	newPerson, err = service.GetById(newPerson.Id)

	if err != nil {
		t.Error("Error getting by id", err)
	}

	checkPerson(t, newPerson, signupPerson)

	newPerson, err = service.Login(signupPerson.EmailAddress, password)

	if err != nil {
		t.Error("Error logging in", err)
	}

	checkPerson(t, newPerson, signupPerson)

	err = service.Delete(newPerson)

	if err != nil {
		t.Error("Error deleting the person", err)
	}

}

func checkPerson(t *testing.T, newPerson *Person, signupPerson *Person) {

	if newPerson.Id == 0 {
		t.Error("Person id should not be null")
	}

	if newPerson.Name != signupPerson.Name {
		t.Error("The name was not saved in the db correctly")
	}

	if newPerson.EmailAddress != signupPerson.EmailAddress {
		t.Error("The email address was not saved in the db correctly")
	}
}

func TestPersonCount(t *testing.T) {

	for i := 1; i <= 10; i++ {
		getTestPerson()
	}

	service, err := GetPersonService()

	if err != nil {
		t.Error("Error creating the person service", err)
	}

	count := service.Count()

	if count == 0 {
		t.Error("There is no people in the system")
	}

}
