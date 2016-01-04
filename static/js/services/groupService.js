// public/js/services/groupService.js
angular.module('groupService', [])

    .factory('Group', function($http) {

        return {
            // get all groups or an individual group
            get: function (id) {
                if (id == null) {
                    return $http.get('../../api/v1/group');
                } else {
                    return $http.get('../../api/v1/group/' + id);
                }
            }
        }

    });
