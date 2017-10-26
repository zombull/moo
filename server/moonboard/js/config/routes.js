/**
 *
 */
moon.config(function($routeProvider, $locationProvider) {
    'use strict';

    $routeProvider.caseInsensitiveMatch = true;
    $routeProvider
    .when('/', {
        templateUrl: 'partials/atc.html',
        controller: 'AirTrafficController as ctrl'
    })
    .when('/p/:grade/:page?', {
        templateUrl: 'partials/peruse.html',
        controller: 'PeruseController as ctrl'
    })
    .when('/t/:grade/:page?', {
        templateUrl: 'partials/peruse.html',
        controller: 'PeruseController as ctrl'
    })
    .when('/s/:setter/:page?', {
        templateUrl: 'partials/setter.html',
        controller: 'SetterController as ctrl'
    })
    .when('/st/:setter/:page?', {
        templateUrl: 'partials/setter.html',
        controller: 'SetterController as ctrl'
    })
    .when('/:problem', {
        templateUrl: 'partials/problem.html',
        controller: 'ProblemController as ctrl'
    })
    .otherwise({
        templateUrl: 'partials/404.html',
        controller: '404Controller as ctrl'
    });

    $locationProvider.html5Mode(true);
});