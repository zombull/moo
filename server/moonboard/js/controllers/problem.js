
/**
 *
 */
moon.controller('ProblemController', function ProblemController($scope, $timeout, $routeParams, $mdDialog, moonboard, database, problems) {
    'use strict';

    problems.reset();

    var shadow = {
        attempts: 0,
        sessions: 0,
    };

    function dialog(event, ctrl, locals, fn) {
        return $mdDialog.show({
            targetEvent: event,
            controller: ctrl + 'Controller',
            controllerAs: 'ctrl',
            locals: locals,
            ariaLabel: ctrl.toLowerCase() + '-dialog',
            templateUrl: '{0}/html/{1}.html'.format(ctrl == 'Confirm' ? 'common' : 'static', ctrl.toLowerCase()),
            clickOutsideToClose: true,
        }).then(function(source) { if (fn) { fn(); }
        }).catch(function() { });
    }

    database.all(function(data) {
        var problemUrl = $routeParams.problem;

        if (problemUrl === 'upgradenowjackass') {
            for (var tkey in data.ticks) {
                database.tick.migrate(tkey, $scope);
            }
            for (var pkey in data.projects) {
                database.project.migrate(pkey, $scope);
            }
            for (var ekey in data.exiles) {
                database.exile.migrate(ekey, $scope);
            }
            return;
        }
        if (!data.problems.hasOwnProperty(problemUrl)) {
            $scope.error = $scope.error || { status: 404, data: 'The problem "' + problemUrl + '" does not exist.' };
            return;
        }

        problems.set(data.index.problems);

        var me = data.problems[problemUrl];
        var problem = data.index.problems[me];
        var setter = data.index.setters[problem.r];
        var grades = data.grades[problem.v / 10];
        var suggested = { setter: [], grade: [] };

        var showProblem = function(i, maxLength) {
            var p = data.index.problems[i];
            return !p.t && !p.e && i != me && (suggested.grade.length + suggested.setter.length) < maxLength;
        };
        _.each(setter.p, function(i) {
            if (showProblem(i, 10)) {
                suggested.setter.push(data.index.problems[i]);
            }
        });
        _.each(grades, function(i) {
            if (showProblem(i, 20)) {
                suggested.grade.push(data.index.problems[i]);
            }
        });
        $scope.setter = setter;
        $scope.problem = problem;
        $scope.suggested = suggested;
        $scope.attempts = shadow.attempts = problem.p ? problem.p.a : 0;
        $scope.sessions = shadow.sessions = problem.p ? problem.p.s : 0;

        moonboard.load().then(
            function() {
                moonboard.set(problem.h);

                if (problem.t && $routeParams.nuke === 'tick') {
                    dialog(undefined, 'Confirm', { prompt: 'Nuke tick?', buttons: { cancel: '2620', confirm: '2622' } }, function() {
                        database.tick.rm(problem.u, $scope);
                    });
                }
            },
            function() {
                $scope.error = $scope.error || { status: 500, data: 'Failed to load Moonboard' };
            }
        );
    }, $scope);

    var updateTimeout = null;
    function queueUpdate() {
        $timeout.cancel(updateTimeout);
        updateTimeout = $timeout(function(problemUrl, attempts, sessions) {
            database.project.add(problemUrl, { a: attempts, s: sessions }, $scope);
            shadow.attempts = attempts;
            shadow.sessions = sessions;
        }, 3000, true, $routeParams.problem, $scope.attempts, $scope.sessions);
    }

    $scope.attempt = function (event) {
        $scope.attempts++;
        if ($scope.sessions === 0) {
            $scope.sessions++;
        }
        queueUpdate();
    };
    $scope.session = function (event) {
        if ($scope.attempts === 0) {
            $scope.attempts++;
        }
        $scope.sessions++;
        queueUpdate();
    };
    $scope.tick = function (event) {
        $timeout.cancel(updateTimeout);
        dialog(event, 'Tick', { problem: $scope.problem, attempts: $scope.attempts || 1, sessions: $scope.sessions || 1 });
    };
    $scope.exile = function (event) {
        dialog(event, 'Confirm', { prompt: 'Exile to the dark side of the moon?', buttons: { cancel: '2694', confirm: '2620' } }, function() {
            database.exile.add($scope.problem.u, $scope);
        });
    };
    $scope.nuke = function (event) {
        if ($scope.problem.e) {
            dialog(event, 'Confirm', { prompt: 'Nuke exile status?', buttons: { cancel: '2620', confirm: '2622' } }, function() {
                database.exile.rm($scope.problem.u, $scope);
                bug.On($scope.problem.e);
            });
            return;
        }
        if (updateTimeout && $timeout.cancel(updateTimeout)) {
            $scope.attempts = shadow.attempts;
            $scope.sessions = shadow.sessions;
        } else {
            dialog(event, 'Confirm', { prompt: 'Nuke attempts and sessions?', buttons: { cancel: '2620', confirm: '2622' } }, function() {
                $scope.attempts = shadow.attempts = 0;
                $scope.sessions = shadow.sessions = 0;
                database.project.rm($scope.problem.u, $scope);
                bug.On($scope.problem.p);
            });
        }
    };
});
