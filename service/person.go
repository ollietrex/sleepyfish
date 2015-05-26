//person.go
package service

import (
	"github.com/ollietrex/sleepyfish/database"
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
)

const saltSize = 16

//The service type for the person
type Person struct {
	Id           int64
	Name         string
	EmailAddress string
}

var (
	//The error returned when someone tries to register an account with an email address that already exists
	DuplicateAccountError = errors.New("An account with that email address already exists")
	//The error returned when some one tried to login with invalid credentials
	InvalidCredentialError = errors.New("The credentials provided are invalid")
)

//Base implementation of the person service interface
type personService struct {
	personRepository database.IPersonRepository
}

//Register a user into the system
func (p *personService) Register(person *Person, password string) (*Person, error) {

	_, err := p.personRepository.GetByEmail(person.EmailAddress)

	if err == nil {
		return nil, DuplicateAccountError
	}

	//Create the salt and password

	passwordSalt, passwordHash := p.generateSaltAndPassword(password)

	dbPerson := &database.Person{
		Name:         person.Name,
		EmailAddress: person.EmailAddress,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}
	//Persist the person
	dbPerson, err = p.personRepository.Persist(dbPerson)

	if err != nil {
		return nil, err
	}

	//Return a new person to the users
	newPerson := &Person{
		Id:           dbPerson.Id,
		Name:         dbPerson.Name,
		EmailAddress: dbPerson.EmailAddress,
	}

	return newPerson, nil
}

///Login a user into the system
func (p *personService) Login(email string, password string) (*Person, error) {
	dbPerson, err := p.personRepository.GetByEmail(email)

	if err != nil {
		fmt.Print(err)
		return nil, InvalidCredentialError
	}

	ok := p.checkPassword(dbPerson.PasswordSalt, dbPerson.PasswordHash, password)

	if !ok {
		return nil, InvalidCredentialError
	}

	newPerson := &Person{
		Id:           dbPerson.Id,
		Name:         dbPerson.Name,
		EmailAddress: dbPerson.EmailAddress,
	}

	return newPerson, nil
}

//Get a user by there unique id
func (p *personService) GetById(personId int64) (*Person, error) {
	dbPerson, err := p.personRepository.GetById(personId)
	if err != nil {
		return nil, err
	}

	newPerson := &Person{
		Id:           dbPerson.Id,
		Name:         dbPerson.Name,
		EmailAddress: dbPerson.EmailAddress,
	}

	return newPerson, nil
}

//Delete a person from the persistent storage
func (p *personService) Delete(person *Person) error {
	dbPerson, err := p.personRepository.GetById(person.Id)
	if err != nil {
		return err
	}
	return p.personRepository.Delete(dbPerson)
}

//Get the count of persons in the system
func (s personService) Count() int64 {
	return s.personRepository.Count()
}

//Generate a salt and hash of the users password, the salt is returned
//first and then the hash
func (p *personService) generateSaltAndPassword(password string) ([]byte, []byte) {
	bytePassword := []byte(password)

	// generate salt from given password
	salt := p.generateSalt(bytePassword)

	// generate password + salt hash to store into database
	combination := string(salt) + string(bytePassword)
	passwordSha := sha512.New()
	io.WriteString(passwordSha, combination)
	hash := passwordSha.Sum(nil)

	return salt, hash
}

//Check the password agaist the salt and hash stored against the user
func (p *personService) checkPassword(salt []byte, hash []byte, password string) bool {
	passwordSalt := salt
	passwordHash := hash
	bytePassword := []byte(password)
	combination := string(passwordSalt) + string(bytePassword)
	passswordSha := sha512.New()
	io.WriteString(passswordSha, combination)
	match := bytes.Equal(passswordSha.Sum(nil), passwordHash)
	return match
}

//Generate a new salt for the user
func (p *personService) generateSalt(secret []byte) []byte {
	buf := make([]byte, saltSize, saltSize+sha512.Size)
	io.ReadFull(rand.Reader, buf)
	hash := sha512.New()
	hash.Write(buf)
	hash.Write(secret)
	return hash.Sum(buf)
}

//The interface definition for the person service methods
type IPersonService interface {
	//Register a person into the system with the specified password
	Register(person *Person, password string) (*Person, error)
	//Login a person into the system with the specified password
	Login(email string, password string) (*Person, error)
	//Get a person by there unique id
	GetById(personId int64) (*Person, error)
	//Delete a person from the persistent storage
	Delete(person *Person) error
	//Provide a count of the amount of people registered in the system
	Count() int64
}

//Get an instance of the person service to deal with the persistence of
//people to the persistent storage
func GetPersonService() (IPersonService, error) {
	if personSrv == nil {
		repo, err := database.GetPersonRepository()
		if err != nil {
			return nil, err
		}
		personSrv = &personService{
			personRepository: repo,
		}
	}
	return personSrv, nil
}
