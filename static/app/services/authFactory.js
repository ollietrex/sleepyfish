(function () {

	'use strict';

	angular.module('sleepfishApp.services.auth',[])

	.factory('AuthFactory', ['$q', '$window', '$http', function ($q, $window, $http) {
	    var storedUserInfo;
		function loginUser(userInfo) {
		    var deferred = $q.defer();
		    var promise = deferred.promise;
		
			storedUserInfo = userInfo
			$window.sessionStorage["userInfo"] = JSON.stringify(userInfo);
				
			$http.defaults.headers.common['X-Session-Id'] = userInfo.accessToken;	
				
			deferred.resolve(200);
            //deferred.reject(status);
		
		    promise.success = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    promise.error = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    return promise;
		}
		function logout() {
			var deferred = $q.defer();	

			$window.sessionStorage["userInfo"] = null;
			userInfo = null;
			deferred.resolve(200);
			//deferred.reject(error);

		    promise.success = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    promise.error = function (fn) {
		        promise.then(fn);
		        return promise;
		    }		
		    return promise;
		}	
		function getUserInfo() {
			
		    var deferred = $q.defer();
		    var promise = deferred.promise;
			
			deferred.resolve(storedUserInfo);
			
		    promise.success = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    promise.error = function (fn) {
		        promise.then(fn);
		        return promise;
		    }
		    return promise;
		}
		return {
			loginUser : loginUser,
			getUserInfo : getUserInfo,
			logout: logout
	    }
	}]);



})();