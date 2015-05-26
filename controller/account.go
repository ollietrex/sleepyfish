//account.go
package controller

import (
	"github.com/ollietrex/sleepyfish/service"
	"net/http"
)

//The json type to login to the system
type Login struct {
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
	RememberMe   bool   `json:"remember_me"`
}

type JsonError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//The login and register response which returns the access token to be used
//in subsequent requests to the system. It also contains the users name
//so this can be display to the user for friendly display
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
}

//The json type to register in the system
type Register struct {
	Name         string `json:"name"`
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
}

//Validate if the login request is valid
func isValidLoginRequest(loginRequest Login) bool {
	if !isValidEmailAddress(loginRequest.EmailAddress) {
		return false
	}
	if !isValidPassword(loginRequest.Password) {
		return false
	}
	return true
}

//Validate if the register request is valid
func isValidRegisterRequest(registerRequest Register) bool {
	if !isValidName(registerRequest.Name) {
		return false
	}
	if !isValidEmailAddress(registerRequest.EmailAddress) {
		return false
	}
	if !isValidPassword(registerRequest.Password) {
		return false
	}
	return true
}

//The account login handler to deal with the authetication of sessions
func AccountLoginHandler(w http.ResponseWriter, r *http.Request) {

	writer := GetRender(w)

	loginRequest := Login{}

	err := paraseJSON(r, &loginRequest)

	if err != nil {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid json"))
		return
	}

	if !isValidLoginRequest(loginRequest) {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid email address or password"))
		return
	}

	personService, err := service.GetPersonService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	person, err := personService.Login(loginRequest.EmailAddress, loginRequest.Password)

	if err != nil {
		if err == service.InvalidCredentialError {
			writer.JSON(http.StatusUnauthorized, getError(http.StatusUnauthorized, "Invalid email address or password"))
		} else {
			writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		}
		return
	}

	sessionService, err := service.GetSessionService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	session, err := sessionService.RegisterSession(person.Id)

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	loginResponse := LoginResponse{
		AccessToken: session.Guid,
		Name:        person.Name,
	}

	writer.JSON(http.StatusOK, loginResponse)
}

//The account register handler to deal with the registration of new accconts
func AccountRegisterHandler(w http.ResponseWriter, r *http.Request) {

	writer := GetRender(w)

	registerRequest := Register{}

	err := paraseJSON(r, &registerRequest)

	if err != nil {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid json"))
		return
	}

	if !isValidRegisterRequest(registerRequest) {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid name, email or password"))
		return
	}

	personService, err := service.GetPersonService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	person := &service.Person{
		Name:         registerRequest.Name,
		EmailAddress: registerRequest.EmailAddress,
	}

	person, err = personService.Register(person, registerRequest.Password)

	if err != nil {
		if err == service.DuplicateAccountError {
			writer.JSON(http.StatusConflict, getError(http.StatusConflict, err.Error()))
		} else {
			writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		}
		return
	}

	sessionService, err := service.GetSessionService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	session, err := sessionService.RegisterSession(person.Id)

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	loginResponse := LoginResponse{
		AccessToken: session.Guid,
		Name:        person.Name,
	}

	writer.JSON(http.StatusOK, loginResponse)

}

func AccountLogoutHandler(w http.ResponseWriter, r *http.Request) {
	writer := GetRender(w)

	session, err := GetSession(r)

	if err != nil {
		writer.JSON(http.StatusUnauthorized, getError(http.StatusUnauthorized, "Invalid session"))
		return
	}

	sessionService, err := service.GetSessionService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	err = sessionService.DeleteSession(session.Guid)

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	writer.JSON(http.StatusOK, getError(http.StatusOK, "You have been logged out"))

}
