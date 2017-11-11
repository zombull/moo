
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

    database.all(function(data) {
        var name = $routeParams.problem;
        if (!data.problems.hasOwnProperty(name)) {
            $scope.error = $scope.error || { status: 404, data: 'The problem "' + name + '" does not exist.' };
            return;
        }

        database.project.get(name, $scope, function(project) {
            shadow.attempts = $scope.attempts = project ? project.a : 0;
            shadow.sessions = $scope.sessions = project ? project.s : 0;
            problems.set(data.index.problems);

            var me = data.problems[name];
            var problem = data.index.problems[me];
            var setter = data.index.setters[problem.e];
            var grades = data.grades[problem.v / 10];
            var suggested = { setter: [], grade: [] };
            _.each(setter.p, function(p) {
                if (p != me && suggested.setter.length < 10 && !data.index.problems[p].t) {
                    suggested.setter.push(data.index.problems[p]);
                }
            });
            _.each(grades, function(p) {
                if (p != me && (suggested.grade.length + suggested.setter.length) < 20 && !data.index.problems[p].t) {
                    suggested.grade.push(data.index.problems[p]);
                }
            });
            $scope.setter = setter;
            $scope.problem = problem;
            $scope.suggested = suggested;

            moonboard.load().then(
                function() {
                    moonboard.set(problem.h);
                },
                function() {
                    $scope.error = $scope.error || { status: 500, data: 'Failed to load Moonboard' };
                }
            );
        });
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
        $mdDialog.show({
            targetEvent: event,
            controller: 'TickController',
            controllerAs: 'ctrl',
            locals: { problem: $scope.problem, attempts: $scope.attempts || 1, sessions: $scope.sessions || 1 },
            ariaLabel: 'tick-dialog',
            templateUrl: 'static/html/tick.html',
            clickOutsideToClose: true,
        });
    };
    $scope.tbd = function (event) {

    };
    $scope.nuke = function (event) {
        if (updateTimeout && $timeout.cancel(updateTimeout)) {
            $scope.attempts = shadow.attempts;
            $scope.sessions = shadow.sessions;
        } else {
            $mdDialog.show({
                targetEvent: event,
                controller: 'ConfirmController',
                controllerAs: 'ctrl',
                locals: { prompt: 'Nuke attempts and sessions?', buttons: { cancel: '2620', confirm: '2694' } },
                ariaLabel: 'confirm-dialog',
                templateUrl: 'common/html/confirm.html',
                clickOutsideToClose: true,
            }).then(function(source) {
                $scope.attempts = shadow.attempts = 0;
                $scope.sessions = shadow.sessions = 0;
                database.project.rm($routeParams.problem, $scope);
            }).catch(function() {});
        }
    };
});
