
/**
 *
 */
moon.controller('ProblemController', function ProblemController($scope, $timeout, $routeParams, $mdDialog, moonboard, database, problems) {
    'use strict';

    problems.reset();

    $scope.attemptCodes = [
        '24FF', '2776', '2777', '2778', '2779', '277A', '277B', '277C', '277D', '277E', '277F',
        '24EB', '24EC', '24ED', '24EE', '24EF', '24F0', '24F1', '24F2', '24F3', '24F4'
    ];

    $scope.sessionCodes = [
        '24EA', '2460', '2461', '2462', '2463', '2464', '2465', '2466', '2467', '2468', '2469',
        '246A', '246B', '246C', '246D', '246E', '246F', '2470', '2471', '2472', '2473',
    ];

    var shadow = {
        attempts: 0,
        sessions: 0,
    };

    database.all(function(data) {
        var name = $routeParams.problem;
        if (!data.p.hasOwnProperty(name)) {
            $scope.error = $scope.error || { status: 404, data: 'The problem "' + name + '" does not exist.' };
            return;
        }

        database.project.get(name, $scope, function(project) {
            shadow.attempts = $scope.attempts = project ? project.attempts : 0;
            shadow.sessions = $scope.sessions = project ? project.sessions : 0;
            problems.set(data.i);

            var me = data.p[name];
            var problem = data.i[me];
            var setter = data.i[problem.e];
            var grades = data.g[problem.v / 10];
            var suggested = { setter: [], grade: [] }
            _.each(setter.p, function(p) {
                if (p != me && suggested.setter.length < 10 && !data.i[p].t) {
                    suggested.setter.push(data.i[p])
                }
            });
            _.each(grades, function(p) {
                if (p != me && (suggested.grade.length + suggested.setter.length) < 20 && !data.i[p].t) {
                    suggested.grade.push(data.i[p]);
                }
            });
            $scope.setter = setter;
            $scope.problem = problem;
            $scope.suggested = suggested;

            moonboard.load().then(
                function() {
                    moonboard.set(problem.h)
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
        updateTimeout = $timeout(function(){
            database.project.set($routeParams.problem, { attempts: $scope.attempts, sessions: $scope.sessions,}, $scope);
            shadow.attempts = $scope.attempts;
            shadow.sessions = $scope.sessions;
        }, 3000);
    }

    $scope.attempt = function (event) {
        $scope.attempts++ || $scope.sessions++;
        queueUpdate();
    };
    $scope.session = function (event) {
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
            templateUrl: 'static/partials/tick.html',
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
            }).catch(function() {});
        }
    };
});
