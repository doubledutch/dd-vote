// public/js/services/snackService.js
angular.module('voteService', [])

    .factory('Vote', function($http) {

        return {

            // get all votes for the current user
            getUserVotes : function(group_name) {
                return $http({
                    method: 'GET',
                    url: '../../api/v1/groups/' + group_name + '/votes/user'
                });
            },

            // save a vote
            save : function(post_id, voteData) {
                return $http({
                    method: 'POST',
                    url: '../../api/v1/posts/' + post_id + '/votes',
                    headers: { 'Content-Type' : 'application/json' },
                    data: voteData
                });
            }
        }
    });
