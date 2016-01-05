// public/js/controllers/adminCtrl.js
angular.module('adminCtrl', ['ngRoute'])

    // inject the Snack service into our controller
    .controller('adminController', function($scope, $http, $location, Admin) {

        $scope.loginData = {};

        var path = $location.path().split("/");
        var groupId = path.pop();
        $scope.loginData.groupId = groupId;
        $scope.groupId = groupId;

        $scope.adminLogin = function() {

            Admin.login($scope.loginData)
                .success(function(data) {
                    if (!data.error) {
                        // reload page now that we've logged in
                        location.reload();
                    } else {
                        $('#error-message').html('Unable to login...</br>' + JSON.stringify(data.message));
                        $('#error-message').show().delay(3000).fadeOut('slow');
                    }
                })
                .error(function(data) {
                    console.log(data);
                });
        }

    })

    .config(function ($routeProvider, $locationProvider) {
        $locationProvider.html5Mode(true);
    });
