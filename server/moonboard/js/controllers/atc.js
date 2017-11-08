/**
 * Air Traffic Controller
 */
moon.controller('AirTrafficController', function AirTrafficController($routeParams, database, problems) {
    'use strict';

    if ($routeParams.password) {
        database.password($routeParams.password);
    }
    database.all(function(data) {
        problems.set(data.index.problems);
    });
});
