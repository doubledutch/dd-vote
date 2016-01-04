// public/js/services/snackService.js
angular.module('snackService', [])

    .factory('Snack', function($http) {

        return {

            // get all snacks for or snacks for an individual group
            get : function(id) {
                if (id == null) {
                    return $http.get('../../api/v1/post');
                } else {
                    return $http({
                        url: '../../api/v1/post',
                        method: "GET",
                        params: {group: id}
                    });
                }
            },

            // save a snack (pass in snack data)
            save : function(snackData) {
                return $http({
                    method: 'POST',
                    url: '../../api/v1/post',
                    headers: { 'Content-Type' : 'application/x-www-form-urlencoded' },
                    data: $.param(snackData)
                });
            }
        }

    });
