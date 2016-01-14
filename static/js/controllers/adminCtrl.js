// public/js/controllers/adminCtrl.js
angular.module('adminCtrl', ['ngRoute'])

    // inject the Snack service into our controller
    .controller('adminController', function($scope, $http, $location, Admin) {

        // show an error message in a box at the top of the page. Set fadeOut to true to disappear after an interval
        $scope.showError = function(text, fadeOut) {
          $('#error-message').html(text);
          $('#error-message').show();
          if (fadeOut) {
            $('#error-message').delay(3000).fadeOut('slow')
          }
        }

        $scope.loginData = {};

        var path = $location.path().split("/");
        var groupId = path.pop();
        $scope.loginData.groupId = groupId;
        $scope.groupId = groupId;

        $scope.adminLogin = function() {
            Admin.login($scope.loginData)
                .success(function(data) {
                    // reload page now that we've logged in
                    location.reload(true);
                })
                .error(function(data) {
                    $scope.showError('Unable to log in...</br>' + data.message, false)
                });
        }

        $scope.adminLogout = function() {

            Admin.logout()
                .success(function(data) {
                  location.reload(true);
                })
                .error(function(data) {
                  console.log(data);
                });
        }

    })

    .config(function ($routeProvider, $locationProvider) {
        // this allows the url path parsing to work
        $locationProvider.html5Mode(true);
    });
