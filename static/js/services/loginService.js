// public/js/services/loginService.js
angular.module('loginService', [])

    .factory('User', function($http) {

        return {
            // attempt to log a user in
            save: function (userData) {
                return $http({
                    method: 'POST',
                    url: '../../prideLogin',
                    headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                    data: $.param(userData)
                });
            }
        }

    });
