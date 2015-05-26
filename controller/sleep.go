package controller

import (
	"github.com/gorilla/mux"
	"github.com/ollietrex/sleepyfish/service"
	"net/http"
	"strconv"
	"time"
)

type Sleep struct {
	Id      int64  `json:"id"`
	Start   int64  `json:"start"`
	End     int64  `json:"end"`
	Quality int8   `json:"quality"`
	Feeling int8   `json:"feeling"`
	Comment string `json:"comment"`
}

func isValidSleep(sleep *Sleep, unixTimeBase int64) bool {

	//End must be after start
	if sleep.Start > sleep.End {
		return false
	}

	//A sleep duration cant be more than 48 hours
	sleepDuration := (sleep.End - sleep.Start) * int64(time.Second)

	if sleepDuration > int64(time.Hour*48) {
		return false
	}

	//We can only update the start and end date 6 months each side of its creation date
	sixMonths := int64(((60 * 60) * 24) * (31 * 6))

	minDate := unixTimeBase - sixMonths
	maxDate := unixTimeBase + sixMonths

	if sleep.Start <= minDate || sleep.Start >= maxDate {
		return false
	}

	if sleep.End <= minDate || sleep.End >= maxDate {
		return false
	}

	if sleep.Quality < 0 || sleep.Quality > 10 {
		return false
	}
	if sleep.Feeling < 0 || sleep.Feeling > 10 {
		return false
	}
	//We ignore comment for the moment
	return true
}

func SleepCreateHandler(w http.ResponseWriter, r *http.Request) {
	writer := GetRender(w)

	session, err := GetSession(r)

	if err != nil {
		writer.JSON(http.StatusUnauthorized, getError(http.StatusUnauthorized, "Invalid session"))
		return
	}

	sleepRequest := &Sleep{}

	err = paraseJSON(r, sleepRequest)

	if err != nil {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid json"))
		return
	}

	if !isValidSleep(sleepRequest, time.Now().Unix()) {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid sleep object"))
		return
	}

	serviceSleep := toServiceSleep(&service.Sleep{}, sleepRequest, session)

	sleepService, err := service.GetSleepService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	serviceSleep, err = sleepService.Persist(serviceSleep)

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	sleepRequest = toRequestSleep(sleepRequest, serviceSleep)

	writer.JSON(http.StatusOK, sleepRequest)
}

func SleepUpdateHandler(w http.ResponseWriter, r *http.Request) {
	writer := GetRender(w)

	session, err := GetSession(r)

	if err != nil {
		writer.JSON(http.StatusUnauthorized, getError(http.StatusUnauthorized, "Invalid session"))
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid id"))
		return
	}

	sleepRequest := &Sleep{}

	err = paraseJSON(r, sleepRequest)

	if err != nil {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid json"))
		return
	}

	sleepService, err := service.GetSleepService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	serviceSleep, err := sleepService.GetById(int64(id), session.PersonId)

	if err != nil {
		writer.JSON(http.StatusNotFound, getError(http.StatusNotFound, "No such sleep in the system"))
		return
	}

	if !isValidSleep(sleepRequest, serviceSleep.CreatedOn) {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid sleep object"))
		return
	}

	serviceSleep = toServiceSleep(serviceSleep, sleepRequest, session)

	serviceSleep, err = sleepService.Persist(serviceSleep)

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	sleepRequest = toRequestSleep(sleepRequest, serviceSleep)

	writer.JSON(http.StatusOK, sleepRequest)

}

func SleepDeleteHandler(w http.ResponseWriter, r *http.Request) {
	writer := GetRender(w)

	session, err := GetSession(r)

	if err != nil {
		writer.JSON(http.StatusUnauthorized, getError(http.StatusUnauthorized, "Invalid session"))
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid id"))
		return
	}

	sleepService, err := service.GetSleepService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	serviceSleep, err := sleepService.GetById(int64(id), session.PersonId)

	if err != nil {
		writer.JSON(http.StatusNotFound, getError(http.StatusNotFound, "No such sleep in the system"))
		return
	}

	err = sleepService.Delete(serviceSleep)

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	writer.JSON(http.StatusOK, getError(http.StatusOK, "Deleted"))

}

func SleepGetHandler(w http.ResponseWriter, r *http.Request) {
	writer := GetRender(w)

	session, err := GetSession(r)

	if err != nil {
		writer.JSON(http.StatusUnauthorized, getError(http.StatusUnauthorized, "Invalid session"))
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid id"))
		return
	}

	sleepService, err := service.GetSleepService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	serviceSleep, err := sleepService.GetById(int64(id), session.PersonId)

	if err != nil {
		writer.JSON(http.StatusNotFound, getError(http.StatusNotFound, "No such sleep in the system"))
		return
	}

	sleepRequest := toRequestSleep(&Sleep{}, serviceSleep)

	writer.JSON(http.StatusOK, sleepRequest)

}

func SleepSearchHandler(w http.ResponseWriter, r *http.Request) {
	writer := GetRender(w)

	err := r.ParseForm()

	if err != nil {
		writer.JSON(http.StatusBadRequest, getError(http.StatusBadRequest, "Invalid query parameters"))
		return
	}

	maxEnd := time.Now().Add(time.Hour * 24 * 365).Unix()

	start := FormValueOrDefault(r, "start", 0)
	end := FormValueOrDefault(r, "end", maxEnd)
	min := FormValueOrDefault(r, "min", 0)
	max := FormValueOrDefault(r, "max", 10000)

	session, err := GetSession(r)

	if err != nil {
		writer.JSON(http.StatusUnauthorized, getError(http.StatusUnauthorized, "Invalid session"))
		return
	}

	sleepService, err := service.GetSleepService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	serviceSleeps, err := sleepService.GetByPerson(session.PersonId, start, end, min, max)

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	sleep := []*Sleep{}

	for _, item := range serviceSleeps {
		newItem := toRequestSleep(&Sleep{}, item)
		sleep = append(sleep, newItem)
	}

	writer.JSON(http.StatusOK, sleep)
}

//Map from the json object to the service object
func toServiceSleep(serviceSleep *service.Sleep, sleepRequest *Sleep, session *service.Session) *service.Sleep {
	serviceSleep.Id = sleepRequest.Id
	serviceSleep.PersonId = session.PersonId
	serviceSleep.Start = sleepRequest.Start
	serviceSleep.End = sleepRequest.End
	serviceSleep.Quality = sleepRequest.Quality
	serviceSleep.Feeling = sleepRequest.Feeling
	serviceSleep.Comment = sleepRequest.Comment
	return serviceSleep
}

//Map from the service object to the json object
func toRequestSleep(sleepRequest *Sleep, serviceSleep *service.Sleep) *Sleep {
	sleepRequest.Id = serviceSleep.Id
	sleepRequest.Start = serviceSleep.Start
	sleepRequest.End = serviceSleep.End
	sleepRequest.Quality = serviceSleep.Quality
	sleepRequest.Feeling = serviceSleep.Feeling
	sleepRequest.Comment = serviceSleep.Comment
	return sleepRequest
}
