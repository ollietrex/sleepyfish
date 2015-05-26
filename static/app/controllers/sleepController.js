(function () {
	
	'use strict';
	
	angular.module('sleepfishApp.controllers.sleep',[])
	
	.controller('SleepIndexCtrl',['$scope', '$state', 'SleepService', function ($scope, $state, SleepService) { 
	
		$scope.data = [];
	
		$scope.retriveSleep = function() {
			SleepService.listSleep().success(function(data) {
				$scope.data = data.sleep;
			}).error(function(status) {
				alert(status)
			});
		};
	
		$scope.retriveSleep();
	
		$scope.deleteSleep = function(sleepId) {
			SleepService.deleteSleep(sleepId).success(function() {
				$scope.retriveSleep();
			}).error(function(status) {
				alert(status);
			});
		};
	
	}])
	.controller('SleepEditCtrl',['$scope', '$state', '$stateParams', 'SleepService', function ($scope, $state, $stateParams, SleepService) {
		$scope.data = {};
		$scope.showError = false; 
		
		$scope.startDateTimeNow = function() {
			$scope.data.start = new Date();
		};
		$scope.endDateTimeNow = function() {
			$scope.data.end = new Date();
		};
  
		$scope.minuteStep = 5;
		$scope.showMeridian = false;

		$scope.startDateTimeNow();
		$scope.endDateTimeNow();
		
		$scope.retriveSleep = function(sleepId) {
			SleepService.getSleep(sleepId).success(function(data) {
				
			$scope.data.start = new Date(data.sleep.start * 1000);
			$scope.data.end = new Date(data.sleep.end * 1000);
			$scope.data.feeling = data.sleep.feeling;
			$scope.data.quality = data.sleep.quality;
			$scope.data.comment = data.sleep.comment;
				
			}).error(function(status) {
				alert(status);
			});
		};

		$scope.edit = function () {
			if (!$scope.sleepForm.$invalid) {
				$scope.showError = false;
				var start = Math.round($scope.data.start.getTime() / 1000);
				var end = Math.round($scope.data.end.getTime() / 1000);
				
				SleepService.updateSleep($stateParams.sleepId, start, end, $scope.data.quality, $scope.data.feeling, $scope.data.comment).success(function (data) {
					$state.go('sleep');
				}).error(function (status) {
					$scope.showError = true;
					if (status == 400 || status == 401 || status == 404) {
						$scope.errorDetails = "Please check you entered values and try again, please try again";					
					} else {
						$scope.errorDetails = "Opps, something went wrong, please try again later";					
					} 
				});
			}
		};	

		$scope.retriveSleep($stateParams.sleepId);
		
	}])
	.controller('SleepCreateCtrl',['$scope', '$state', 'SleepService', function ($scope, $state, SleepService) { 
		$scope.data = {};
		$scope.showError = false;
		$scope.errorDetails = "";		
		
		$scope.startDateTimeNow = function() {
			$scope.data.start = new Date();
		};
		$scope.endDateTimeNow = function() {
			$scope.data.end = new Date();
		};
  
		$scope.minuteStep = 5;
		$scope.showMeridian = false;

		$scope.startDateTimeNow();
		$scope.endDateTimeNow();

		$scope.create = function () {
			if (!$scope.sleepForm.$invalid) {
				$scope.showError = false;
				var start = Math.round($scope.data.start.getTime() / 1000);
				var end = Math.round($scope.data.end.getTime() / 1000);
				
				SleepService.createSleep(start, end, $scope.data.quality, $scope.data.feeling, $scope.data.comment).success(function (data) {
					$state.go('sleep');
				}).error(function (status) {
					$scope.showError = true;
					if (status == 400 || status == 401 || status == 404) {
						$scope.errorDetails = "Please check you entered values and try again, please try again";					
					} else {
						$scope.errorDetails = "Opps, something went wrong, please try again later";					
					} 
				});
			}
		};	
	}]);

})();