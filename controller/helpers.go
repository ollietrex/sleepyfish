package controller

import (
	"github.com/ollietrex/sleepyfish/service"
	"errors"
	"net/http"
	"strconv"
)

func FormValueOrDefault(r *http.Request, value string, defaultValue int64) int64 {
	valueStr := r.FormValue(value)
	if len(valueStr) == 0 {
		return defaultValue
	}
	valueInt, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return valueInt
}

//We private two methods to get the session
func GetSession(r *http.Request) (*service.Session, error) {
	notFoundError := errors.New("Invalid session")

	sessionId := r.Header.Get("X-Session-Id")

	if len(sessionId) == 0 {
		cookie, err := r.Cookie("s")
		if err == nil && len(cookie.Value) == 0 {
			sessionId = cookie.Value
		}
	}

	if len(sessionId) == 0 {
		return nil, notFoundError
	}

	sessionService, err := service.GetSessionService()

	if err != nil {
		return nil, err
	}

	session, err := sessionService.GetSession(sessionId)

	if err != nil {
		return nil, err
	}

	return session, nil
}
