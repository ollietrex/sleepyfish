package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

//Parase the json into an object
func paraseJSON(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)

	if err != nil {
		return err
	}

	return nil
}

//Return an error object to the user which formats in a standard format
//for the application
func getError(status int, message string) JsonError {
	return JsonError{
		Status:  status,
		Message: message,
	}
}

//
func isValidEmailAddress(emailAddress string) bool {
	if len(emailAddress) == 0 {
		return false
	}
	atIndex := strings.LastIndex(emailAddress, "@")

	if atIndex == -1 {
		return false
	}

	dotIndex := strings.LastIndex(emailAddress, ".")

	if dotIndex == -1 {
		return false
	}

	//If the last dot is before the at then we know there is an error
	//in the email address
	if dotIndex <= atIndex {
		return false
	}

	return true
}

func isValidPassword(password string) bool {
	if len(password) <= 4 {
		return false
	}
	return true
}

func isValidName(name string) bool {
	if len(name) == 0 {
		return false
	}
	return true
}
