// public/js/services/snackService.js
angular.module('snackService', [])

    .factory('Snack', function($http) {

        return {

            // get all snacks for or snacks for an individual group
            get : function(group_name) {
                return $http({
                    url: '../../api/v1/groups/' + group_name + '/posts',
                    method: "GET"
                });
            },

            // save a snack (pass in snack data)
            save : function(snackData, group_name) {
                return $http({
                    method: 'POST',
                    url: '../../api/v1/groups/' + group_name + '/posts',
                    headers: { 'Content-Type' : 'application/json' },
                    data: snackData
                });
            },

            remove: function(snackId) {
              return $http({
                  url: '../../api/v1/posts/' + snackId,
                  method: "DELETE"
              });
            }
        }

    });
