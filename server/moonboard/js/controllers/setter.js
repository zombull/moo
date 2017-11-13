
/**
 *
 */
moon.controller('SetterController', function SetterController($scope, $location, $routeParams, bug, browse, database) {
    'use strict';

    var __problems = []; // Local list used as the source for problems.
    var showTicks = ($location.path().split('/')[1].toLowerCase() === 'st');

    $scope.title = showTicks ? 'Ticks' : 'Problems';

    if (!browse.ready($scope, $routeParams.page, 0)) {
        return;
    }

    database.all(function(data) {
        var skey = 's/' + $routeParams.setter.toLowerCase();
        if (!data.setters.hasOwnProperty(skey)) {
            $scope.error = $scope.error || { status: 404, data: 'Did not find a setter matching "' + $routeParams.setter + '"' };
            return;
        }

        var setter = data.index.setters[data.setters[skey]];
        $scope.setter = _.pick(setter, ['u', 'n']);
        if (!showTicks) {
            $scope.setter.u = 'st/' + $routeParams.setter.toLowerCase();
        }

        _.each(setter.p, function(i) {
            if (!showTicks == !data.index.problems[i].t) {
                __problems.push(data.index.problems[i]);
            }
        });
        if (__problems.length === 0) {
            bug.on(!showTicks);
            $scope.error = $scope.error || { status: 404, data: 'Did not find any ticked problems set by {0}.'.format(setter.n) };
            return;
        }

        browse.go($scope, __problems, function() {});
    });
});
