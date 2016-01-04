// public/js/services/adminService.js
angular.module('adminService', [])

    .factory('Admin', function($http) {

        return {
            login: function (userData) {
                return $http({
                    method: 'POST',
                    url: '../../admin/login',
                    headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                    data: $.param(userData)
                });
            }
        }

    });
