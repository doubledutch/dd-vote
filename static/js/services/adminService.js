// public/js/services/adminService.js
angular.module('adminService', [])

    .factory('Admin', function($http) {

        return {
            login: function (userData) {
                return $http({
                    method: 'POST',
                    url: '../../api/v1/admin/login',
                    headers: {'Content-Type': 'application/json'},
                    data: userData
                });
            }
        }

    });
