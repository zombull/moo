
/**
 *
 */
moon.controller('SetterController', function SetterController($scope, $location, $routeParams, bug, browse, database, history) {
    'use strict';

    var __problems = []; // Local list used as the source for problems.
    var rp = $location.path().split('/')[1].toLowerCase();
    var showTicks = (rp === 'st');
    var historyKey = rp + '.' + $routeParams.setter.toLowerCase();

    $scope.title = showTicks ? 'Ticks' : 'Problems';

    database.all(function(data) {
        var skey = 's/' + $routeParams.setter.toLowerCase();
        if (!data.setters.hasOwnProperty(skey)) {
            $scope.error = $scope.error || { status: 404, data: 'Did not find a setter matching "' + $routeParams.setter + '"' };
            return;
        }

        if (!browse.ready($scope, $routeParams.page, historyKey)) {
            return;
        }

        var setter = data.index.setters[data.setters[skey]];
        $scope.setter = _.pick(setter, ['u', 'n']);
        if (!showTicks) {
            $scope.setter.u = 'st/' + $routeParams.setter.toLowerCase();
        }

        _.each(setter.p, function(i) {
            var p = data.index.problems[i];
            if (!p.e && !showTicks == !p.t) {
                __problems.push(p);
            }
        });
        if (__problems.length === 0) {
            bug.on(!showTicks);
            $scope.error = $scope.error || { status: 404, data: 'Did not find any ticked problems set by {0}.'.format(setter.n) };
            return;
        }

        browse.go($scope, __problems, function(i) {
            history.set(historyKey, $scope.i, true);
        });
    });
});
