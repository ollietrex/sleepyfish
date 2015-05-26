(function () {

	'use strict';
	
	var app = angular.module('sleepyfishApp', [
		'ui.router',
		'ui.bootstrap', 
		'ui.bootstrap.datetimepicker',
		'sleepfishApp.services.auth',
		'sleepfishApp.services.account',
		'sleepfishApp.services.sleep',
		'sleepfishApp.controllers.home',
		'sleepfishApp.controllers.account',
		'sleepfishApp.controllers.sleep',
		'sleepfishApp.directives'
		])
	
	.config(['$stateProvider', '$urlRouterProvider', function ($stateProvider,   $urlRouterProvider) {
	
		var authenticated = ['$q', 'AuthFactory', function ($q, AuthFactory) {
			var deferred = $q.defer();
			
			var userInfo = AuthFactory.getUserInfo().success(function(userInfo) {
				if (userInfo != null) {
					deferred.resolve();
				} else {
					deferred.reject('Not logged in');
				}
			}).error(function() {
				deferred.reject('Not logged in');				
			})

			return deferred.promise;
		}];
	
	
	
		$urlRouterProvider.otherwise('/home');
	
		$stateProvider
	
		    .state('home', {
		        url: '/home',
		        templateUrl: '/static/app/views/home/index.html',
		        controller: 'HomeIndexCtrl'
		    })
		
		    .state('login', {
		        url: '/account/login',
		        templateUrl: '/static/app/views/account/login.html',
		        controller: 'AccountLoginCtrl'
		    })
		
		    .state('register', {
		        url: '/account/register',
		        templateUrl: '/static/app/views/account/register.html',
		        controller: 'AccountRegisterCtrl'
		    })
	
		    .state('sleep', {
		        url: '/sleep',
		        templateUrl: '/static/app/views/sleep/index.html',
		        controller: 'SleepIndexCtrl',
				resolve: {
					authenticated: authenticated
				}
		    })
		    .state('sleepCreate', {
		        url: '/sleep/create',
		        templateUrl: '/static/app/views/sleep/create.html',
		        controller: 'SleepCreateCtrl',
				resolve: {
					authenticated: authenticated
				}
		    })
		    .state('sleepEdit', {
		        url: '/sleep/edit/{sleepId:int}',
		        templateUrl: '/static/app/views/sleep/edit.html',
		        controller: 'SleepEditCtrl',
				resolve: {
					authenticated: authenticated
				}
		    });
	}])
	
	.run(["$rootScope", "$state", function($rootScope, $state) {
		$rootScope.$on('$stateChangeError', function (err, req) {
			$state.go('login');
		});
	}]);

})();
