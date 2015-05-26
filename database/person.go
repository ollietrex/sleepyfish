//person.go
package database

import (
	"gopkg.in/gorp.v1"
	"time"
)

//The data poco for a person
type Person struct {
	Id           int64  `db:"id"`
	Name         string `db:"name"`
	EmailAddress string `db:"email"`
	PasswordHash []byte `db:"password_hash"`
	PasswordSalt []byte `db:"password_salt"`
	UpdatedOn    int64  `db:"updated_on"`
	CreatedOn    int64  `db:"created_on"`
}

//Pre insert hook to updated on and created on date for an insert
func (i *Person) PreInsert(s gorp.SqlExecutor) error {
	i.UpdatedOn = time.Now().Unix()
	i.CreatedOn = i.UpdatedOn
	return nil
}

//Pre update hook to update the updated date for an update
func (i *Person) PreUpdate(s gorp.SqlExecutor) error {
	i.UpdatedOn = time.Now().Unix()
	return nil
}

//The base implementation of the person repository
type personRepository struct {
	db *gorp.DbMap
}

//Get the person by there unique if from the repository
func (p personRepository) GetById(id int64) (*Person, error) {
	obj, err := p.db.Get(Person{}, id)
	if err != nil {
		return nil, err
	}
	person := obj.(*Person)
	return person, nil
}

//Get a person from the persistent storage by there email address
func (p personRepository) GetByEmail(email string) (*Person, error) {
	person := &Person{}
	err := p.db.SelectOne(&person, "select * from person where email=:email", map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}
	return person, nil
}

//Persist a person to the persistent respository
func (p personRepository) Persist(person *Person) (*Person, error) {
	var err error
	if person.Id == 0 {
		err = p.db.Insert(person)
	} else {
		_, err = p.db.Update(person)
	}
	if err != nil {
		return nil, err
	}
	return person, nil
}

//Delete a person from the persistent storage
func (p personRepository) Delete(person *Person) error {
	_, err := p.db.Delete(person)
	return err
}

//Count of the amount of people registered in the system
func (s personRepository) Count() int64 {
	count, _ := s.db.SelectInt("SELECT COUNT(*) FROM person")
	return count
}

//The interface definition for the persistence of persons
//into the persistent storage
type IPersonRepository interface {
	//Get a person from the persistent storage by there unique identifier
	GetById(id int64) (*Person, error)
	//Get a person from the persistent storage by there email address
	GetByEmail(email string) (*Person, error)
	//Persist a person to the persistent storage
	Persist(person *Person) (*Person, error)
	//|Delete a person from the persistent storage
	Delete(person *Person) error
	//Count of the amount of people registered in the system
	Count() int64
}

//Get an instance of the person repository for to deal with the
//persistentce of people into the persistent storage
func GetPersonRepository() (IPersonRepository, error) {
	if personRepo == nil {
		db, err := GetDatabase()
		if err != nil {
			return nil, err
		}
		personRepo = &personRepository{
			db: db,
		}
	}
	return personRepo, nil
}
