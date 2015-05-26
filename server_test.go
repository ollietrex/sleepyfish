package main

import (
	"github.com/ollietrex/sleepyfish/controller"
	"github.com/ollietrex/sleepyfish/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	server *httptest.Server
	reader io.Reader
)

//Setup the http server to accept new connections
func init() {
	server = httptest.NewServer(Handlers())
}

//Test the home page renders correctly
func TestHome(t *testing.T) {
	status, _, err := ProcessHttp("/", "GET", "")
	//Ensure the home page returns ok
	CheckResult(t, http.StatusOK, status, err)
}

//Test the registration and login of a user into the system
func TestRegisterLoginSuccess(t *testing.T) {
	email := helpers.RandomEmail()
	password := helpers.RandomString(8)
	registerJson := fmt.Sprintf(`{"name": "%s", "email": "%s", "password": "%s" }`, helpers.RandomName(), email, password)
	loginJson := fmt.Sprintf(`{"email": "%s", "password": "%s", "remember_me" : true}`, email, password)
	//Create the account
	AccountPost(t, registerJson, http.StatusOK)
	//Try and login with the account
	AccountLoginPost(t, loginJson, http.StatusOK)
}

//Test the login functionality
func TestLoginValidation(t *testing.T) {
	AccountLoginPost(t, ``, http.StatusBadRequest)
	AccountLoginPost(t, `{`, http.StatusBadRequest)
	AccountLoginPost(t, `{}`, http.StatusBadRequest)
	AccountLoginPost(t, `{"email": "ollie@tribe.guru" }`, http.StatusBadRequest)
	AccountLoginPost(t, `{"password": "Password13!"}`, http.StatusBadRequest)
	AccountLoginPost(t, `{"remember_me" : true}`, http.StatusBadRequest)
	AccountLoginPost(t, `{"email": "", "password": "Password13!", "remember_me" : true}`, http.StatusBadRequest)
	AccountLoginPost(t, `{"email": "ollie@tribe.guru", "password": "", "remember_me" : true}`, http.StatusBadRequest)
}

//The account login functionality to test posts to the the server
func AccountLoginPost(t *testing.T, json string, expectedStatus int) {
	status, body, err := ProcessHttp(ApiAuthAuthenticate, "POST", json)
	t.Logf("Login: %s", body)
	CheckResult(t, expectedStatus, status, err)
}

//Test the register functionality
func TestRegister(t *testing.T) {

	registerJson := fmt.Sprintf(`{"name": "%s", "email": "%s", "password": "%s" }`, helpers.RandomName(), helpers.RandomEmail(), helpers.RandomString(8))

	AccountPost(t, registerJson, http.StatusOK)

	AccountPost(t, registerJson, http.StatusConflict)

	AccountPost(t, ``, http.StatusBadRequest)
	AccountPost(t, `{`, http.StatusBadRequest)
	AccountPost(t, `{}`, http.StatusBadRequest)
	AccountPost(t, `{"name": "" }`, http.StatusBadRequest)

	//No data
	AccountPost(t, `{"name": "", "email": "", "password": "" }`, http.StatusBadRequest)
	//No email
	AccountPost(t, `{"name": "Joe", "email": "", "password": "william" }`, http.StatusBadRequest)
	//No password
	AccountPost(t, `{"name": "Joe Blogs", "email": "joe@blogs.com", "password": "" }`, http.StatusBadRequest)
	//Password to short
	AccountPost(t, `{"name": "Joe Blogs", "email": "joe@blogs.com", "password": "wi" }`, http.StatusBadRequest)
	//No @ in password
	AccountPost(t, `{"name": "Joe Blogs", "email": "joeblogs.com", "password": "william" }`, http.StatusBadRequest)
	//No dot in password
	AccountPost(t, `{"name": "Joe Blogs", "email": "joe@blogscom", "password": "william" }`, http.StatusBadRequest)
	//Dot before @
	AccountPost(t, `{"name": "Joe Blogs", "email": "joe.blogs@com", "password": "william" }`, http.StatusBadRequest)

}

//The generic method on the to post data to the account register for testing
func AccountPost(t *testing.T, json string, expectedStatus int) {
	status, body, err := ProcessHttp(ApiAuthRegister, "POST", json)

	t.Logf("Register: %s", body)

	CheckResult(t, expectedStatus, status, err)
}

//Test the sleep life cycle to ensure all functionality works
func TestSleep(t *testing.T) {

	loginResponse, err := GetLoginForTest()

	if err != nil {
		t.Error(err)
	}

	createSleep, err := CreateSleepForTest(loginResponse)

	if err != nil {
		t.Error(err)
	}

	if createSleep.Id == 0 {
		t.Error("Did not create sleep")
	}

	getSleep, err := GetSleepForTest(loginResponse, createSleep.Id)

	if err != nil {
		t.Error(err)
	}

	if getSleep.Id == 0 {
		t.Error("Could not get sleep")
	}

	createSleep.Comment = "New comment"

	updateSleep, err := UpdateSleepForTest(loginResponse, createSleep)

	if err != nil {
		t.Error(err)
	}

	if createSleep.Comment != updateSleep.Comment {
		t.Error("Did not update the comment correctly")
	}

	err = DeleteSleepForTest(loginResponse, getSleep.Id)

	if err != nil {
		t.Error(err)
	}

	_, err = GetSleepForTest(loginResponse, createSleep.Id)

	if err == nil {
		t.Error("Could not delete the sleep")
	}

}

//Make a request to delete sleep from the server
func DeleteSleepForTest(loginResponse *controller.LoginResponse, id int64) error {
	url := strings.Replace(ApiSleepsWithId, "{id:[0-9]+}", strconv.Itoa(int(id)), -1)

	status, _, err := ProcessHttpWithSession(url, "DELETE", loginResponse.AccessToken, "")

	if status != http.StatusOK {
		return errors.New(fmt.Sprintf("Http Status %s", status))
	}

	if err != nil {
		return err
	}

	return nil
}

//Get sleep from the server as a test to ensure that it works
func GetSleepForTest(loginResponse *controller.LoginResponse, id int64) (*controller.Sleep, error) {
	url := strings.Replace(ApiSleepsWithId, "{id:[0-9]+}", strconv.Itoa(int(id)), -1)

	status, body, err := ProcessHttpWithSession(url, "GET", loginResponse.AccessToken, "")

	if status != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Http Status %s", status))
	}

	if err != nil {
		return nil, err
	}

	controllerSleep := &controller.Sleep{}

	err = json.Unmarshal([]byte(body), controllerSleep)

	if err != nil {
		return nil, err
	}

	return controllerSleep, nil

}

//Make a request to create sleep on the server
func CreateSleepForTest(loginResponse *controller.LoginResponse) (*controller.Sleep, error) {

	sleepJson := fmt.Sprintf(`{"start": %d, "end": %d, "quality":1, "feeling":1, "comment":"This is a comment"}`, time.Now().Unix(), time.Now().Unix())

	status, body, err := ProcessHttpWithSession(ApiSleeps, "POST", loginResponse.AccessToken, sleepJson)

	if status != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Http Status %d", status))
	}

	if err != nil {
		return nil, err
	}

	controllerSleep := &controller.Sleep{}

	err = json.Unmarshal([]byte(body), controllerSleep)

	if err != nil {
		return nil, err
	}

	return controllerSleep, nil
}

//Make a request to update sleep on the server to ensure it works
func UpdateSleepForTest(loginResponse *controller.LoginResponse, sleep *controller.Sleep) (*controller.Sleep, error) {

	url := strings.Replace(ApiSleepsWithId, "{id:[0-9]+}", strconv.Itoa(int(sleep.Id)), -1)

	sleepJson, err := json.Marshal(sleep)

	fmt.Println(string(sleepJson))
	fmt.Println("{\"id\":91,\"start\":1432018524,\"end\":1432018524,\"quality\":4,\"feeling\":1,\"comment\":\"Test\"}")

	status, body, err := ProcessHttpWithSession(url, "PUT", loginResponse.AccessToken, string(sleepJson))

	if status != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Http Status %d", status))
	}

	if err != nil {
		return nil, err
	}

	controllerSleep := &controller.Sleep{}

	err = json.Unmarshal([]byte(body), controllerSleep)

	if err != nil {
		return nil, err
	}

	return controllerSleep, nil
}

//Make a request to the server for login so it can be retuened to the user
func GetLoginForTest() (*controller.LoginResponse, error) {

	email := helpers.RandomEmail()
	password := helpers.RandomString(8)

	registerJson := fmt.Sprintf(`{"name": "%s", "email": "%s", "password": "%s" }`, helpers.RandomName(), email, password)
	loginJson := fmt.Sprintf(`{"email": "%s", "password": "%s", "remember_me" : true}`, email, password)

	status, _, err := ProcessHttp(ApiAuthRegister, "POST", registerJson)

	if status != http.StatusOK {
		return nil, errors.New("Cound not create user for login")
	}

	status, body, err := ProcessHttp(ApiAuthAuthenticate, "POST", loginJson)

	if status != http.StatusOK {
		return nil, errors.New("Cound not login")
	}

	if err != nil {
		return nil, err
	}

	loginResponse := &controller.LoginResponse{}

	err = json.Unmarshal([]byte(body), loginResponse)

	if err != nil {
		return nil, err
	}
	return loginResponse, nil
}

func TestSleepTest(t *testing.T) {
	loginResponse, err := GetLoginForTest()

	for i := 0; i < 20; i++ {
		CreateSleepForTest(loginResponse)
	}

	status, body, err := ProcessHttpWithSession(ApiSleeps, "GET", loginResponse.AccessToken, "")

	CheckResult(t, http.StatusOK, status, err)

	sleep := []controller.Sleep{}

	err = json.Unmarshal([]byte(body), &sleep)

	if err != nil {
		t.Error(err)
	}

	if len(sleep) == 0 {
		t.Error("Could not load any sleep")
	}
}

//Make a http request to the server
func ProcessHttp(url string, method string, json string) (int, string, error) {
	return ProcessHttpWithSession(url, method, "", json)
}

//Make a http request to the server with a session added as a header
func ProcessHttpWithSession(url string, method string, session string, json string) (int, string, error) {

	fullUrl := fmt.Sprintf("%s%s", server.URL, url)

	var request *http.Request

	if len(json) == 0 {
		request, _ = http.NewRequest(method, fullUrl, nil)
	} else {
		reader = strings.NewReader(json)
		request, _ = http.NewRequest(method, fullUrl, reader)
	}

	if len(session) > 0 {
		request.Header.Set("X-Session-Id", session)
	}

	res, err := http.DefaultClient.Do(request)

	data, _ := ioutil.ReadAll(res.Body)

	return res.StatusCode, string(data), err

}

//Check the result against the expected result and output an error if one is shown
func CheckResult(t *testing.T, expectedStatus, status int, err error) {
	if status != expectedStatus {
		if err != nil {
			t.Errorf("Failed with status: %d, expected %d with error: %s", status, expectedStatus, err.Error())
		} else {
			t.Errorf("Failed with status: %d, expected %d", status, expectedStatus)
		}
	}
}
