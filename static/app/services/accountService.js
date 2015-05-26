(function () {
	
	'use strict';
	
	angular.module('sleepfishApp.services.account',[])

	.factory('AccountService', ['$q','$http', function ($q, $http) {
		function loginUser(name, password) {
		    var deferred = $q.defer();
		    var promise = deferred.promise;
		
		    $http.post("/api/v1/auth/authenticate", { "email": name, "password": password, "remember_me": true }).
		        success(function(data, status, headers, config) {
					var userInfo = {
							accessToken: data.access_token,
							userName: data.name
					};
					deferred.resolve({
						status: status,
						userInfo: userInfo
					});
		        }).
		        error(function(data, status, headers, config) {
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
		function registerUser(name, email, password) {
		    var deferred = $q.defer();
		    var promise = deferred.promise;
		
		    $http.post("/api/v1/auth/register", { "name": name, "email": email, "password": password }).
		        success(function (data, status, headers, config) {
					var userInfo = {
							accessToken: data.access_token,
							userName: data.name
					};
					deferred.resolve({
						status: status,
						userInfo: userInfo
					});
		        }).
		        error(function (data, status, headers, config) {
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
		function logout(accessToken) {
			var deferred = $q.defer();	
			$http({
				method: "POST",
				url: logoutUrl,
				headers: {
					"X-SessionId": accessToken
				}
			}).then(function(result) {
				deferred.resolve(result);
			}, function(error) {
				deferred.reject(error);
			});
		
			return deferred.promise;
		}	
		return {
			loginUser : loginUser,
			registerUser : registerUser,
			logout: logout
	    }
	}]);

})();