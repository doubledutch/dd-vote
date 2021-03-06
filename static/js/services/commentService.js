// public/js/services/commentService.js
angular.module('commentService', [])

    .factory('Comment', function($http) {

        return {
            // save a comment
            save: function (snackId, commentData) {
                return $http({
                    method: 'POST',
                    url: '../../api/v1/posts/' + snackId + '/comments',
                    headers: {'Content-Type': 'application/json'},
                    data: commentData
                });
            }
        }

    });
