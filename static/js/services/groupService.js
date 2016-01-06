// public/js/services/groupService.js
angular.module('groupService', [])

    .factory('Group', function($http) {

        return {
            // get an individual group
            get: function (id) {
              return $http.get('../../api/v1/groups/' + id);
            }
        }

    });
