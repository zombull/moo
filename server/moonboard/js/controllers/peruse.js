
/**
 *
 */
moon.controller('PeruseController', function PeruseController($scope, $location, $routeParams, bug, browse, database, history) {
    'use strict';

    var __problems = []; // Local list used as the source for problems.
    var rp = $location.path().split('/')[1].toLowerCase();
    var showTicks = (rp === 't');
    var showProjects = (rp === 'j');
    var historyKey = rp + '.' + $routeParams.grade.toLowerCase();

    $scope.title = showProjects ? 'Projects' : showTicks ? 'Ticks' : 'Problems';

    if (!browse.ready($scope, $routeParams.page, historyKey)) {
        return;
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

    function problemsFromMap(m, data) {
        var p = [];
        _.each(m, function(v, k) {
            bug.on(!data.problems.hasOwnProperty(k));
            var problem = data.index.problems[data.problems[k]];
            if (!grade || problem.g === grade) {
                p.push(problem);
            }
        });
        return p;
    }

    database.all(function(data) {
        if (showProjects) {
            __problems = problemsFromMap(data.projects, data);
        } else if (showTicks) {
            __problems = problemsFromMap(data.ticks, data);
            __problems.sort(function(a, b) {
                return (new Date(b.t.d)) - (new Date(a.t.d));
            });
        } else {
            __problems = data.index.problems.filter(function(problem) {
                return  (!grade || problem.g === grade) && (!showTicks == !problem.t);
            });
        }

        if (__problems.length === 0) {
            var meta = $routeParams.grade === 'all' ? '' : $routeParams.grade + ' ';
            $scope.error = $scope.error || { status: 404, data: 'Did not find any {0}{1}.'.format(meta, $scope.title) };
            return;
        }

        browse.go($scope, __problems, function(i) {
            $scope.setter = data.index.setters[$scope.problem.e];
            history.set(historyKey, $scope.i);
        });
    }, $scope);
});
