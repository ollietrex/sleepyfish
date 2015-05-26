//count.go
package controller

import (
	"bitbucket.org/ollietrex/sleepyfish/service"
	"net/http"
)

//The json type to login to the system
type Count struct {
	PeopleCount int64 `json:"peopleCount"`
	SleepCount  int64 `json:"sleepCount"`
}

func CountIndexHandler(w http.ResponseWriter, r *http.Request) {

	writer := GetRender(w)

	personService, err := service.GetPersonService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	sleepService, err := service.GetSleepService()

	if err != nil {
		writer.JSON(http.StatusInternalServerError, getError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	count := Count{
		PeopleCount: personService.Count(),
		SleepCount:  sleepService.Count(),
	}

	writer.JSON(http.StatusOK, count)

}
