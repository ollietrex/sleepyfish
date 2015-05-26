package views

import (
	"github.com/ollietrex/sleepyfish/models"
	"bytes"
)

func HomeIndex(model *models.HomeIndexViewModel) string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n<!DOCTYPE html>\n\n<!--[if lt IE 7]>      <html lang=\"en\" ng-app=\"sleepyfishApp\" class=\"no-js lt-ie9 lt-ie8 lt-ie7\"> <![endif]-->\n\n<!--[if IE 7]>         <html lang=\"en\" ng-app=\"sleepyfishApp\" class=\"no-js lt-ie9 lt-ie8\"> <![endif]-->\n\n<!--[if IE 8]>         <html lang=\"en\" ng-app=\"sleepyfishApp\" class=\"no-js lt-ie9\"> <![endif]-->\n\n<!--[if gt IE 8]><!--> <html lang=\"en\" ng-app=\"sleepyfishApp\" class=\"no-js\"> <!--<![endif]-->\n\n<head lang=\"en\">\n\n    <meta charset=\"utf-8\">\n\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n\n    <title>Sleep Fish</title>\n\n	\n\n	<link rel=\"stylesheet\" href=\"/content/css\">\n\n\n\n    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->\n\n    <!--[if lt IE 9]>\n\n      <script src=\"//oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js\"></script>\n\n      <script src=\"//oss.maxcdn.com/respond/1.4.2/respond.min.js\"></script>\n\n    <![endif]-->\n\n	\n\n</head>\n\n<body>\n\n	<header class=\"navbar navbar-default navbar-fixed-top bs-docs-nav\" id=\"top\" role=\"banner\">\n\n		<div class=\"container\">\n\n			<div class=\"navbar-header\">\n\n				<button class=\"navbar-toggle collapsed\" type=\"button\" data-toggle=\"collapse\" data-target=\".bs-navbar-collapse\">\n\n					<span class=\"sr-only\">Toggle navigation</span>\n\n					<span class=\"icon-bar\"></span>\n\n					<span class=\"icon-bar\"></span>\n\n					<span class=\"icon-bar\"></span>\n\n				</button>\n\n				<a href=\"#/home\" class=\"navbar-brand\">Sleep Fish</a>\n\n			</div>\n\n			<nav class=\"collapse navbar-collapse bs-navbar-collapse\">\n\n				<ul class=\"nav navbar-nav\">\n\n					<li>\n\n						<a href=\"#/sleep\">Sleep</a>\n\n					</li>\n\n				</ul>\n\n				<ul class=\"nav navbar-nav navbar-right\">\n\n					<li><a href=\"#/account/login\">Login</a></li>\n\n					<li><a href=\"#/account/register\">Register</a></li>\n\n				</ul>\n\n			</nav>\n\n		</div>\n\n	</header>\n\n    <div class=\"container\">\n\n	  <div ui-view></div>\n\n    </div>\n\n	<script src=\"/content/js\"></script>\n\n</body>\n\n</html>")

	return _buffer.String()
}
