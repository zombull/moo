/**
 *
 */
moon.config(function($routeProvider, $locationProvider) {
    'use strict';

    $routeProvider.caseInsensitiveMatch = true;
    $routeProvider
    .when('/', {
        templateUrl: 'static/html/atc.html',
        controller: 'AirTrafficController as ctrl'
    })
    .when('/p/:grade', {
        templateUrl: 'static/html/peruse.html',
        controller: 'PeruseController as ctrl'
    })
    .when('/t/:grade', {
        templateUrl: 'static/html/peruse.html',
        controller: 'PeruseController as ctrl'
    })
    .when('/j/:grade', {
        templateUrl: 'static/html/peruse.html',
        controller: 'PeruseController as ctrl'
    })
    .when('/s/:setter', {
        templateUrl: 'static/html/peruse.html',
        controller: 'SetterController as ctrl'
    })
    .when('/st/:setter', {
        templateUrl: 'static/html/peruse.html',
        controller: 'SetterController as ctrl'
    })
    .when('/:problem', {
        templateUrl: 'static/html/problem.html',
        controller: 'ProblemController as ctrl'
    })
    .otherwise({
        templateUrl: 'static/html/404.html',
        controller: '404Controller as ctrl'
    });

    $locationProvider.html5Mode(true);
});