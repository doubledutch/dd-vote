// public/js/services/commentService.js
angular.module('commentService', [])

    .factory('Comment', function($http) {

        return {
            // save a comment
            save: function (commentData) {
                return $http({
                    method: 'POST',
                    url: '../../api/v1/comment',
                    headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                    data: $.param(commentData)
                });
            }
        }

    });
