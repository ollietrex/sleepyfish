(function () {

	'use strict';
	
	angular.module('sleepfishApp.controllers.account',[])
    
	.controller('AccountLoginCtrl',['$scope', '$state','AuthFactory','AccountService', function ($scope, $state, AuthFactory, AccountService) { 
		$scope.data = {};
		$scope.showError = false;
		$scope.errorDetails = "";
		$scope.login = function () {
			if (!$scope.loginForm.$invalid) {
			
				AccountService.loginUser($scope.data.email, $scope.data.password).success(function (data) {
					
					//One the have logged the user in we register them with the auth factory
					AuthFactory.loginUser(data.userInfo).success(function(state) {
						$state.go('sleep');						
					}).error(function(state) {
						$scope.showError = true;
						$scope.errorDetails = "Opps, something went wrong, please try again later";											
					})
					
				}).error(function (status) {
					$scope.showError = true;
					if (status == 400 || status == 401 || status == 404) {
						$scope.errorDetails = "Incorrect email address or password, please try again";					
					} else {
						$scope.errorDetails = "Opps, something went wrong, please try again later";					
					} 
				});			
			}
		};
	
	}])
	
	.controller('AccountRegisterCtrl', ['$scope', '$state', 'AuthFactory', 'AccountService', function ($scope, $state, AuthFactory, AccountService) { 
		$scope.data = {};
		$scope.showError = false;
		$scope.errorDetails = "";	
		$scope.register = function () {
			if (!$scope.registerForm.$invalid) {
				
				AccountService.registerUser($scope.data.name, $scope.data.email, $scope.data.password).success(function (data) {
					
					//One the have logged the user in we register them with the auth factory
					AuthFactory.loginUser(data.userInfo).success(function(state) {
						$state.go('sleep');						
					}).error(function(state) {
						$scope.showError = true;
						$scope.errorDetails = "Opps, something went wrong, please try again later";											
					})
					
				}).error(function (status) {
					$scope.showError = true;
					if (status == 401 || status == 404) {
						$scope.errorDetails = "Please check your details and try again";					
					} else if (status == 409) {
						$scope.errorDetails = "An account with this email address is already registered";										
					} else {
						$scope.errorDetails = "Opps, something went wrong, please try again later";					
					} 
				});		
			}
		};
	
	}]);
	

})();