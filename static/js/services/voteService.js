// public/js/services/snackService.js
angular.module('voteService', [])

    .factory('Vote', function($http) {

        return {

            // get all votes for the current user
            getUserVotes : function() {
                return $http({
                    method: 'GET',
                    url: '../../api/v1/user/votes'
                });
            },

            // save a vote
            save : function(voteData) {
                return $http({
                    method: 'POST',
                    url: '../../api/v1/vote',
                    headers: { 'Content-Type' : 'application/json' },
                    data: voteData
                });
            }
        }
    });
