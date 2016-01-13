// public/js/controllers/mainCtrl.js
angular.module('mainCtrl', ['ngRoute'])

    // inject the Snack service into our controller
    .controller('mainController', function($scope, $http, $location, $anchorScroll, $routeParams, $timeout, $interval, Snack, User, Group, Vote) {

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

        // show an error message in a box at the top of the page. Set fadeOut to true to disappear after an interval
        $scope.showError = function(text, fadeOut) {
          $('#error-message').html(text);
          $('#error-message').show();
          if (fadeOut) {
            $('#error-message').delay(3000).fadeOut('slow')
          }
        }

        // put a post into Presentation Mode
        $scope.present = function(snack) {
            $scope.presentedSnack = snack;
            $scope.isPresenting = true;
        }

        // queue up new data and show the user the "Tap to load new data" floating box
        $scope.queueData = function(data) {
            $('#new-data-toast').stop().fadeIn(400);
            $scope.queuedData = data.value;
        }

        // load queued data into the current dataset
        $scope.loadQueuedData = function() {
            $('#new-data-toast').stop().fadeOut(400);
            $scope.snacks = $scope.queuedData;
        }

        // set default sort values for the tabs
        $scope.predicate = 'sum_votes';
        $scope.reverse = true;

        // object to hold selected group
        $scope.selected = {};

        // object to hold all the data for the new snack form
        $scope.snackData = {};
        $scope.voteData = {};

        // loading variable to show the spinning loading icon
        $scope.loading = true;
        $scope.votesLoading = true;

        // get group id from path
        var path = $location.path().split("/");
        var groupId = path.pop();
        var pathType = path.pop();
        $scope.moderation = $location.search().moderation === 'true';

        // queued data to load
        $scope.queuedData = [];

        // get all the questions in a group first and bind it to the $scope.snacks
        var loadQuestions = function() {
            Snack.get(groupId)
                .success(function(data) {
                    for (var i in data.value) {
                        // calculate sum of votes
                        data.value[i].sum_votes = data.value[i].upvotes - data.value[i].downvotes;
                    }
                    if ($scope.loading) {
                        // we were showing the loading spinner, now show the data
                        $scope.snacks = data.value;
                        $scope.loading = false;
                    } else if ($scope.autorefresh) {
                        // auto-refresh mode is turned on, so automatically update data
                        $scope.snacks = data.value;
                    } else if (!angular.equals($scope.snacks, data.value)) {
                        // queue up new data and show appropriate message
                        $scope.queueData(data);
                    }
                });
        };

        // get all of the user's votes for questions in this group and put them
        // in a hashmap mapping post UUIDs to vote values
        var loadUserVotes = function() {
          // only load vote data once
          if ($scope.votesLoading) {
            Vote.getUserVotes(groupId)
                .success(function(data) {
                    $scope.votesLoading = false;
                    for (var i in data.value) {
                        // update our hashmap of post UUIDs with the user's vote value
                        $scope.voteData[data.value[i].post_uuid] = data.value[i].value;
                    }
                });
          }
        }

        // get questions data and user data for group
        var loadData = function() {
            loadQuestions();
            loadUserVotes();
        }

        // load data and continue to get new data every 15s
        var loadDataWithInterval = function() {
          loadData();

          // continue to get new data on an interval
          $interval(function () {
              loadData();
          }, 15000);
        }

        // assume js sdk is initialized, but show error after .5 seconds if
        // it doesn't actually get initialized in that time
        $scope.isJSSDKInitialized = true;
        var showJSSDKError = $timeout(function () {
            $scope.isJSSDKInitialized = false;
            loadDataWithInterval();
        }, 500);
        DD.Events.onReady(function() {
            // mark js sdk as initialized
            $timeout.cancel(showJSSDKError);
            $scope.isJSSDKInitialized = true;
            DD.Events.getCurrentUserAsync(function (user) {
                var newUser = new Object();
                newUser.userId = user.UserId || user.Id; // Android passes userId and iOS passes Id
                newUser.userId = parseInt(newUser.userId); // Android passes userId as a string
                newUser.firstName = user.FirstName;
                newUser.lastName = user.LastName;
                User.save(newUser)
                    .success(function (data) {
                        loadDataWithInterval();
                    })
                    .error(function(data) {
                        $scope.showError('Unable to log in...</br>' + data.message, false)
                    });
            });
        });

        // function to handle submitting the form
        // SAVE A SNACK ======================================================
        $scope.submitSnack = function() {

            var MAX_LENGTH = 140;

            // check if question is too long
            if ($scope.snackData.name.length > MAX_LENGTH) {
                $scope.showError('Questions cannot be longer than ' + MAX_LENGTH + ' characters', true)
                return;
            }

            // save the snack. pass in snack data from the form
            // use the function we created in our service
            Snack.save($scope.snackData, groupId)
                .success(function(data) {
                    // add snack to list
                    snack = data.value;
                    snack.sum_votes = snack.upvotes - snack.downvotes;
                    snack.comments = [];
                    $scope.snacks.push(snack);

                    // clear input
                    $scope.snackData.name = null;
                })
                .error(function(data) {
                    $scope.showError('Unable to submit question, please try again...</br>' + data.message, true)
                });
        }

        $scope.deleteSnack = function(snack) {
          // confirm that the user wants to delete the post
          bootbox.confirm('Delete "' + snack.name + '"?', function(result) {
            if (!result) {
              return;
            }

            Snack.remove(snack.uuid)
                .success(function(data) {
                  // remove the approriate post
                  $.each($scope.snacks, function(i) {
                      if($scope.snacks[i].uuid === snack.uuid) {
                          $scope.snacks.splice(i,1);
                          return false;
                      }
                  });
                })
                .error(function(data) {
                  $scope.showError('Unable to delete question: ' + data.message, true)
                });
          });
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
                if ($scope.snacks[index].uuid == snackId) {

                    // remove old vote value
                    if ($scope.voteData[snackId] == -1) {
                        $scope.snacks[index].downvotes--;
                    } else if ($scope.voteData[snackId] == 1) {
                        $scope.snacks[index].upvotes--;
                    }

                    // add new vote
                    if (value == -1) {
                        $scope.snacks[index].downvotes++;
                    } else if (value == 1) {
                        $scope.snacks[index].upvotes++;
                    }

                    // set new vote value
                    $scope.voteData[snackId] = value;
                    $scope.snacks[index].sum_votes = $scope.snacks[index].upvotes - $scope.snacks[index].downvotes;
                    break;
                }
            }

            Vote.save(snackId, {value: value})
                .success(function(data) {
                    // do nothing because we have already updated our local data
                })
                .error(function(data) {
                    // we could rollback here, but it's probably not worth it
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
            Comment.save(snackId, {comment: $scope.commentText})
                .success(function (data) {
                    // add to snack's comments
                    for (index = 0; index < $scope.snacks.length; index++) {
                        if ($scope.snacks[index].uuid == snackId) {
                            $scope.snacks[index].comments.push(data.value);
                            break;
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
        // this allows the url path parsing to work
        $locationProvider.html5Mode(true);
    });
