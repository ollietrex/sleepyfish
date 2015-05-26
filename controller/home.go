//home.go
package controller

import (
	"bitbucket.org/ollietrex/sleepyfish/models"
	"bitbucket.org/ollietrex/sleepyfish/views"
	"net/http"
)

func HomeIndex(w http.ResponseWriter, r *http.Request) {
	writer := GetRender(w)

	viewModel := &models.HomeIndexViewModel{
		Name: "Sleepy Fish",
	}

	writer.HTML(http.StatusOK, views.HomeIndex(viewModel))
}
