// public/js/controllers/mainCtrl.js
angular.module('mainCtrl', ['ngRoute'])

    // inject the Snack service into our controller
    .controller('mainController', function($scope, $http, $location, $anchorScroll, $routeParams, $timeout, $interval, Snack, User, Group) {

        // sort by controvery rating
        $scope.Math = window.Math;
        $scope.controversySort = function(snack) {
            if (snack.upvotes <= 0 || snack.downvotes <=0) {
                return 0;
            }
            // maximize total votes relative to difference in points, and adjust slightly for total votes
            // this seems better for a smaller number of votes
            total = parseInt(snack.upvotes) + parseInt(snack.downvotes);
            return ((total + 1) / (Math.abs(snack.upvotes - snack.downvotes) + 1)) * Math.sqrt(total);
        };
        // alternative Reddit algorithm
        $scope.redditControversySort = function(snack) {
            if (snack.upvotes <= 0 || snack.downvotes <=0) {
                return 0;
            }
            // maximize for number of votes and a close ratio between the number of upvotes and downvotes
            return (parseInt(snack.upvotes) + parseInt(snack.downvotes)) / (Math.max(snack.upvotes, snack.downvotes) / Math.min(snack.upvotes, snack.downvotes));
        }

        $scope.timeSince = function timeSince(date) {
            return moment.min(moment.utc(date), moment()).fromNow();
        };

        // set default sort values
        $scope.predicate = 'sum_votes';
        $scope.reverse = true;

        // object to hold selected group
        $scope.selected = {};

        // object to hold all the data for the new snack form
        $scope.snackData = {};

        // loading variable to show the spinning loading icon
        $scope.loading = true;

        // get group id from path
        var path = $location.path().split("/");
        var groupId = path.pop();
        var pathType = path.pop();

        // queued data to load
        $scope.queuedData = [];

        var loadData = function() {
            // get all the snacks first and bind it to the $scope.snacks object
            // use the function we created in our service
            // GET ALL SNACKS ====================================================
            Snack.get(groupId)
                .success(function(data) {
                    // calculate sum of votes
                    for (var i in data.snacks) {
                        data.snacks[i].sum_votes = data.snacks[i].upvotes - data.snacks[i].downvotes;
                    }
                    if ($scope.loading) {
                        $scope.snacks = data.snacks;
                        $scope.loading = false;
                    }
                    else if (!angular.equals($scope.snacks, data.snacks)) {
                        $('#new-data-toast').stop().fadeIn(400);
                        $scope.queuedData = data.snacks;
                    }
                });
        };

        // assume js sdk is initialized, but show error after .5 seconds if
        // it doesn't actually get initialized in that time
        $scope.isJSSDKInitialized = true;
        var showJSSDKError = $timeout(function () {
            $scope.isJSSDKInitialized = false;
            loadData();
        }, 500);
        DD.Events.onReady(function() {
            // mark js sdk as initialized
            $timeout.cancel(showJSSDKError);
            $scope.isJSSDKInitialized = true;
            DD.Events.getCurrentUserAsync(function (user) {
                var newUser = new Object();
                newUser.userId = user.UserId || user.Id;
                newUser.firstName = user.FirstName;
                newUser.lastName = user.LastName;
                User.save(newUser)
                    .success(function (data) {
                        // check for failure
                        if (data.error) {
                            $('#error-message').show();
                            $('#error-message').html('Unable to log in...</br>' + JSON.stringify(data.message));
                            return;
                        }
                        loadData();

                        // continue to get new data on an interval
                        $interval(function () {
                            loadData();
                        }, 15000);
                    })
            });
        });

        // function to handle submitting the form
        // SAVE A SNACK ======================================================
        $scope.submitSnack = function() {

            var MAX_LENGTH = 140;

            // check if question is too long
            if ($scope.snackData.name.length > MAX_LENGTH) {
                $('#error-message').html('Questions cannot be longer than ' + MAX_LENGTH + ' characters');
                $('#error-message').show().delay(3000).fadeOut('slow');
                return;
            }
        
            // set group id
            $scope.snackData.group_name = groupId;
            var origValue = $scope.snackData.name;

            // save the snack. pass in snack data from the form
            // use the function we created in our service
            Snack.save($scope.snackData)
                .success(function(data) {
                    if (!data.error) {
                        // add snack to list
                        snack = data.snack;
                        snack.sum_votes = snack.upvotes - snack.downvotes;
                        snack.vote_value = 0;
                        snack.comments = [];
                        $scope.snacks.push(snack);
                    } else {
                        $('#error-message').html('Unable to submit question, please try again...</br>' + JSON.stringify(data.message));
                        $('#error-message').show().delay(3000).fadeOut('slow');
                        $scope.snackData.name = origValue;
                    }

                    // scroll snack into view
                    //$location.hash('snack-' + snack.id);
                    //$anchorScroll();
                })
                .error(function(data) {
                    console.log(data);
                });


            // clear input
            $scope.snackData.name = null;
        }

        $scope.present = function(snack) {
            $scope.presentedSnack = snack;
            $scope.isPresenting = true;
        }

        $scope.loadQueuedData = function() {
            $('#new-data-toast').stop().fadeOut(400);
            $scope.snacks = $scope.queuedData;
        }
    })
    // inject the Vote service into our controller
    .controller('voteController', function($scope, $http, Vote) {

        $scope.vote = function(snackId, value) {

            if (!$scope.isJSSDKInitialized) {
                alert('You can only vote within the DoubleDutch app');
                return;
            }

            // update votes values so the user gets instant feedback
            // this means we make the user think their vote worked immediately
            for (index = 0; index < $scope.snacks.length; index++ ){
                if ($scope.snacks[index].id == snackId) {

                    // remove old vote value
                    if ($scope.snacks[index].vote_value == -1) {
                        $scope.snacks[index].downvotes--;
                    } else if ($scope.snacks[index].vote_value == 1) {
                        $scope.snacks[index].upvotes--;
                    }

                    // add new vote
                    if (value == -1) {
                        $scope.snacks[index].downvotes++;
                    } else if (value == 1) {
                        $scope.snacks[index].upvotes++;
                    }

                    // set new vote value
                    $scope.snacks[index].vote_value = value;
                    $scope.snacks[index].sum_votes = $scope.snacks[index].upvotes - $scope.snacks[index].downvotes;
                    break;
                }
            }

            Vote.save({id: snackId, value: value})
                .success(function(data) {
                    // do nothing
                })
                .error(function(data) {
                    console.log(data);
                });
        }

    })

    // inject the Group service into our controller
    .controller('groupController', function($scope, $http, $location, Group) {
        // make architecture better
    })

    // inject the comment service into our controller
    .controller('commentController', function($scope, $http, Comment) {

        $scope.submitComment = function (snackId) {

            // save the snack. pass in snack data from the form
            // use the function we created in our service
            Comment.save({id: snackId, comment: $scope.commentText})
                .success(function (data) {
                    if (!data.error) {
                        // add to snack's comments
                        for (index = 0; index < $scope.snacks.length; index++) {
                            if ($scope.snacks[index].id == data.comment.snack_id) {
                                $scope.snacks[index].comments.push(data.comment);
                                break;
                            }
                        }
                    }
                })
                .error(function (data) {
                    console.log(data);
                });


            // clear input
            $scope.commentText = null;
        }
    })

    .config(function ($routeProvider, $locationProvider) {
        $routeProvider
            .when('/index.php/g/:groupId', {
                controller: 'MainController'
            });

        // configure html5 to get links working on jsfiddle
        $locationProvider.html5Mode(true);


    });