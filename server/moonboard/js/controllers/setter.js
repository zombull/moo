
/**
 *
 */
moon.controller('SetterController', function SetterController($scope, $location, $routeParams, moonboard, database, problems) {
    'use strict';

    problems.reset();

    $scope.problem = null;
    $scope.i = 0; // Current index into __problems

    var __problems = []; // Local list used as the source for problems.
    var perpage = 15;
    var showTicks = ($location.path().split('/')[1].toLowerCase() === 'st');

    if ($routeParams.page) {
        var page = parseInt($routeParams.page);
        if (isNaN(page)) {
            $scope.error = $scope.error || { status: 404, data: '"' + $routeParams.page + '" is not a page number.' };
            return;
        }
        $scope.i = $routeParams.page * perpage;
    }

    database.all(function(data) {
        var skey = 's/' + $routeParams.setter.toLowerCase();
        if (!data.s.hasOwnProperty(skey)) {
            $scope.error = $scope.error || { status: 404, data: 'Did not find a setter matching "' + $routeParams.setter + '"' };
            return;
        }

        $scope.setter = data.i[data.s[skey]];

        _.each($scope.setter.p, function(i) {
            if (!showTicks == !data.i[i].t) {
                __problems.push(data.i[i]);
            }
        });
        if (__problems.length === 0) {
            bug.on(!showTicks);
            $scope.error = $scope.error || { status: 404, data: 'Did not find any ticked problems set by ' + $routeParams.setter + '.' };
            return;
        }
        problems.set(__problems);
        $scope.count = __problems.length;

        moonboard.load().then(
            function() {
                update(Math.min($scope.i, __problems.length - 1));
            },
            function() {
                $scope.error = $scope.error || { status: 500, da2ta: 'Failed to load Moonboard' };
            }
        );
    }, $scope);

    function update(i) {
        $scope.i = i;
        $scope.problem = __problems[$scope.i];
        moonboard.set($scope.problem.h);

        if (__problems.length > perpage) {
            $scope.list = [];
            var start = Math.min($scope.i, __problems.length - perpage - 1);
            $scope.list = _.slice(__problems, start, start + perpage);
        } else {
            $scope.list = __problems;
        }
    }

    $scope.ppage = function (event) {
        update(Math.max($scope.i - perpage, 0));
    };
    $scope.prev = function (event) {
        if ($scope.i > 0) {
            update($scope.i-1);
        }
    };
    $scope.rand = function (event) {
        update(Math.floor(Math.random() * __problems.length));
    };
    $scope.next = function (event) {
        if ($scope.i < (__problems.length - 1)) {
            update($scope.i+1);
        }
    };
    $scope.npage = function (event) {
        update(Math.min($scope.i + perpage, __problems.length - 1));
    };
});
