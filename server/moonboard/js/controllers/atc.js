/**
 * Air Traffic Controller
 */
moon.controller('AirTrafficController', function AirTrafficController(database, problems) {
    'use strict';

    database.all(function(data) {
        problems.set(data.index.problems);
    });
});
