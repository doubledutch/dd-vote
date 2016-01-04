// public/js/services/commentService.js
angular.module('commentService', [])

    .factory('Comment', function($http) {

        return {
            // save a comment
            save: function (snackId, commentData) {
                return $http({
                    method: 'POST',
                    url: '../../api/v1/comment',
                    headers: {'Content-Type': 'application/json'},
                    data: commentData,
                    params: {post: snackId}
                });
            }
        }

    });
