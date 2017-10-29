
/**
 *
 */
moon.controller('PeruseController', function PeruseController($scope, $location, $routeParams, moonboard, database, problems, bug) {
    'use strict';

    problems.reset();

    $scope.problem = null;
    $scope.i = 0; // Current index into __problems

    var __data = {}; // The global data list, needed to retrieve setter info.
    var __problems = []; // Local list used as the source for problems.
    var perpage = 15;
    var rp = $location.path().split('/')[1].toLowerCase();
    var showTicks = (rp === 't' || rp === 'k');
    var showTocks = (rp === 'k');
    var showProjects = (rp === 'j');

    if ($routeParams.page) {
        var page = parseInt($routeParams.page);
        if (isNaN(page)) {
            $scope.error = $scope.error || { status: 404, data: '"' + $routeParams.page + '" is not a page number.' };
            return;
        }
        $scope.i = $routeParams.page * perpage;
    }

    var grade = $routeParams.grade.toUpperCase();
    if (grade === 'ALL') {
        grade = false;
    } else {
        var vgrade = parseInt(grade.substring(1));
        if (grade.substring(0, 1) !== 'V' || isNaN(vgrade) || vgrade < 4 || vgrade > 17) {
            $scope.error = $scope.error || { status: 404, data: '"' + $routeParams.grade + '" is not a valid grade: must be V4-V17 or ALL.' };
            return;
        }
    }

    database.all(function(data) {
        __data = data;

        if (showProjects) {
            _.each(data.projects, function(v, k) {
                bug.on(!data.problems.hasOwnProperty(k));
                var problem = data.index[data.problems[k]];
                if (!grade || problem.g === grade) {
                    __problems.push(problem);
                }
            });
        } else {
            __problems = _.slice(data.index, 0, _.size(data.problems)).filter(function(problem) {
                return  (!grade || problem.g === grade) && (!showTicks == !problem.t) && (!showTocks || problem.t.hasOwnProperty('p'));
            });
        }

        if (__problems.length === 0) {
            var meta = $routeParams.grade === 'all' ? '' : $routeParams.grade + ' ';
            var type = showProjects ? 'projects' : showTocks ? 'tocks' : showTicks ? 'ticks' : 'problems';
            $scope.error = $scope.error || { status: 404, data: 'Did not find any {0}{1}.'.format(meta, type) };
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
        $scope.setter = __data.index[$scope.problem.e];

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
