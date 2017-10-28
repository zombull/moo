moon.controller('TickController', function($scope, $mdDialog, $filter, database, problem, attempts, sessions) {
    $scope.tick = {
        problem: problem.n,
        grade: problem.g,
        date: new Date(),
        stars: undefined,
        attempts: attempts,
        sessions: sessions,
    };

    $scope.stars = [ 1, 2, 3, 5 ];
    $scope.grades = [ 'V3', 'V4', 'V5', 'V6', 'V7', 'V8', 'V9', 'V10', 'V11', 'V12', 'V13', 'V14', 'V15', 'V16', 'V17' ];

    $scope.cancel = function() {
        $mdDialog.cancel();
    };

    $scope.tock = function() {
        // The variable names/keys matter as they need to match the actual
        // tick schema so that a "tock" can be used interchangeably with a
        database.tock.add({
            p: problem.u,
            d: $filter('date')($scope.tick.date, 'LLLL dd, yyyy'),
            g: $scope.tick.grade,
            s: parseInt($scope.tick.stars),
            a: $scope.tick.attempts,
            e: $scope.tick.sessions > 1 ? $scope.tick.sessions : undefined,
        }, $scope, function() { $mdDialog.hide(); });
    };
});
