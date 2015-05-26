package controller

import (
	"io/ioutil"
	"net/http"
	"os"
)

var (
	coreJsOutput  = ""
	appJsOutput   = ""
	coreCssOutput = ""
)

//The method to serve all of the javascript combined into a single file
func ContentJsHandler(w http.ResponseWriter, r *http.Request) {
	writer := GetRender(w)

	coreJs := []string{}

	if IsProduction() {
		coreJs = append(coreJs, "static/components/jquery/dist/jquery.min.js")
		coreJs = append(coreJs, "static/components/bootstrap/dist/js/bootstrap.min.js")
		coreJs = append(coreJs, "static/components/angular/angular.min.js")
		coreJs = append(coreJs, "static/components/angular-ui-router/release/angular-ui-router.min.js")
		coreJs = append(coreJs, "static/components/angular-bootstrap/ui-bootstrap-tpls.min.js")
		coreJs = append(coreJs, "static/components/angular-ui-bootstrap-datetimepicker/datetimepicker.js")
	} else {
		coreJs = append(coreJs, "static/components/jquery/dist/jquery.js")
		coreJs = append(coreJs, "static/components/bootstrap/dist/js/bootstrap.js")
		coreJs = append(coreJs, "static/components/angular/angular.js")
		coreJs = append(coreJs, "static/components/angular-ui-router/release/angular-ui-router.js")
		coreJs = append(coreJs, "static/components/angular-bootstrap/ui-bootstrap-tpls.js")
		coreJs = append(coreJs, "static/components/angular-ui-bootstrap-datetimepicker/datetimepicker.js")
	}

	if IsProduction() {
		if len(coreJsOutput) == 0 {
			coreJsOutput = getFileContent(coreJs)
		}
	} else {
		coreJsOutput = getFileContent(coreJs)
	}

	appJs := []string{
		"static/components/ie10-viewport-bug-workaround.js",
		"static/app/directives/compareToDirective.js",
		"static/app/services/accountService.js",
		"static/app/services/authFactory.js",
		"static/app/services/sleepService.js",
		"static/app/controllers/homeController.js",
		"static/app/controllers/accountController.js",
		"static/app/controllers/sleepController.js",
		"static/app/app.js",
	}

	if IsProduction() {
		if len(appJsOutput) == 0 {
			appJsOutput = getFileContent(appJs)
		}
	} else {
		appJsOutput = getFileContent(appJs)
	}
	writer.Content(http.StatusOK, "application/javascript", coreJsOutput+"\n\r"+appJsOutput)
}

func ContentCssHandler(w http.ResponseWriter, r *http.Request) {
	writer := GetRender(w)
	css := []string{
		"static/components/bootstrap/dist/css/bootstrap.min.css",
		"static/components/bootstrap/dist/css/bootstrap-theme.min.css",
		"static/components/bootstrap/dist/css/bootstrap-theme.min.css",
		"static/components/angular-ui-bootstrap-datetimepicker/datetimepicker.css",
		"static/css/site.css",
	}

	if IsProduction() {
		if len(coreCssOutput) == 0 {
			coreCssOutput = getFileContent(css)
		}
		writer.Content(http.StatusOK, "text/css", coreCssOutput)
	} else {
		writer.Content(http.StatusOK, "text/css", getFileContent(css))
	}

}

func IsProduction() bool {
	return os.Getenv("PLATFORM") == "production"
}

func getFileContent(files []string) string {
	output := ""
	for _, path := range files {
		data, _ := ioutil.ReadFile(path)
		output += "\n\r" + string(data)
	}
	return output
}
