// public/js/services/loginService.js
angular.module('loginService', [])

    .factory('User', function($http) {

        return {
            // attempt to log a user in
            save: function (userData) {
                return $http({
                    method: 'POST',
                    url: '../../api/v1/login',
                    headers: {'Content-Type': 'application/json'},
                    data: userData
                });
            }
        }

    });
