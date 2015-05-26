(function () {

	'use strict';


	angular.module('sleepfishApp.services.sleep',[])

	.factory('SleepService', ['$q','$http', function ($q, $http) {

		function createSleep(start, end, quality, feeling, comment) {
		    var deferred = $q.defer();
		    var promise = deferred.promise;
			
			var sleep = { 
				"start": start, 
				"end": end, 
				"quality": quality,
				"feeling": feeling,
				"comment": comment
			};
			
			$http.post("/api/v1/sleeps", sleep).success(function (data, status, headers, config) {
				deferred.resolve({
					status: status,
					sleep: data
				});
			}).error(function (data, status, headers, config) {
				deferred.reject(status);
			});

		    promise.success = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    promise.error = function (fn) {
		        promise.then(null, fn);
		        return promise;
		    }
		    return promise;
		}
		function updateSleep(sleepId, start, end, quality, feeling, comment) {
		    var deferred = $q.defer();
		    var promise = deferred.promise;

			var sleep = { 
				"id": sleepId,
				"start": start, 
				"end": end, 
				"quality": quality,
				"feeling": feeling,
				"comment": comment
			};
			
			$http.put("/api/v1/sleeps/" + sleepId, sleep).success(function (data, status, headers, config) {
				deferred.resolve({
					status: status,
					sleep: data
				});
			}).error(function (data, status, headers, config) {
				deferred.reject(status);
			});

		    promise.success = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    promise.error = function (fn) {
		        promise.then(null, fn);
		        return promise;
		    }
		    return promise;
		}
		function deleteSleep(sleepId) {
		    var deferred = $q.defer();
		    var promise = deferred.promise;

			$http.delete("/api/v1/sleeps/" + sleepId).success(function (data, status, headers, config) {
				deferred.resolve(status);
			}).error(function (data, status, headers, config) {
				deferred.reject(status);
			});
			
		    promise.success = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    promise.error = function (fn) {
		        promise.then(null, fn);
		        return promise;
		    }
		    return promise;
		}
		function getSleep(sleepId) {
		    var deferred = $q.defer();
		    var promise = deferred.promise;
			
			$http.get("/api/v1/sleeps/" + sleepId).success(function (data, status, headers, config) {
				deferred.resolve({
					status: status,
					sleep: data
				});
			}).error(function (data, status, headers, config) {
				deferred.reject(status);
			});
			
		    promise.success = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    promise.error = function (fn) {
		        promise.then(null, fn);
		        return promise;
		    }
		    return promise;			
		}
		function listSleep() {
		    var deferred = $q.defer();
		    var promise = deferred.promise;
			
			$http.get("/api/v1/sleeps").success(function (data, status, headers, config) {
				deferred.resolve({
					status: status,
					sleep: data
				});
			}).error(function (data, status, headers, config) {
				deferred.reject(status);
			});
			
		    promise.success = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    promise.error = function (fn) {
		        promise.then(null, fn);
		        return promise;
		    }
		    return promise;
		}
		return {
			createSleep : createSleep,
			updateSleep : updateSleep,
			deleteSleep: deleteSleep,
			listSleep : listSleep,
			getSleep: getSleep
	    }
	}]);


})();