
/**
 *
 */
moon.controller('QueryController', function QueryController($scope, $location, $routeParams, bug, browse, database, history, inspector, problems) {
    'use strict';

    var __problems = []; // Local list used as the source for problems.
    var historyKey = 'q.' + $routeParams.query;
    $scope.title = 'Problems';

    if (!$routeParams.query) {
        $scope.error = $scope.error || { status: 400, data: 'Must specify query for custom browsing' };
        return;
    }

    if (!browse.ready($scope, null, historyKey)) {
        return;
    }

    var query = decodeURI($routeParams.query);
    database.all(function(data) {
        inspector.search(query, true).then(function(results) {
            __problems = results;

            if (__problems.length === 0) {
                $scope.error = $scope.error || { status: 404, data: "Did not find any problems matching query '{0}'".format(query) };
                return;
            }

            browse.go($scope, __problems, function(i) {
                $scope.setter = data.index.setters[$scope.problem.e];
                history.set(historyKey, $scope.i, true);
            });
        }, function() { });
    });
});
