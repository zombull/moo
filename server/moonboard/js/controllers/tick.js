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
        // The variable names/keys matter as the object will be directly
        // serialized to storage, i.e. this is the tick schema.  The date
        // form matters as it's the shortest storage that angular's date
        // filter will accept as input (dates are displayed with a more
        // user friendly input, just cutting down on bytes here...).
        database.tick.add(problem.u, {
            d: $filter('date')($scope.tick.date, 'yyyy-MM-dd'),
            g: $scope.tick.grade,
            s: parseInt($scope.tick.stars),
            a: $scope.tick.attempts,
            e: $scope.tick.sessions > 1 ? $scope.tick.sessions : undefined,
        }, $scope, function() { $mdDialog.hide(); });
    };
});
