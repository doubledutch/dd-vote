// public/js/app.js
var snackApp = angular.module('snackApp', ['mainCtrl', 'snackService', 'voteService', 'groupService',
    'commentService', 'loginService', 'ngSanitize']);

var adminApp = angular.module('adminApp', ['adminCtrl', 'adminService']);

$(function() {
    $('#snack-input').maxlength();
});