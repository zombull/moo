/**
 *
 */
moon.controller('404Controller', function FourOhFourController($scope, $location, database) {
    'use strict';

    $scope.error = {
        status: 404,
        statusText: 'Not Found',
        data: '"' + $location.url() + '" is not a valid URL.'
    };
});