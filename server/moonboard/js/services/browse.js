moon.factory('browse', function (moonboard, problems) {
    'use strict';

    var perpage = 15;

    return {
        ready: function(scope, startPage, alt) {
            problems.reset();

            scope.problem = null;

            if (startPage) {
                var page = parseInt(startPage);
                if (isNaN(page)) {
                    scope.error = scope.error || { status: 404, data: '"' + startPage + '" is not a page number.' };
                    return false;
                }
                scope.i = startPage * perpage;
            } else {
                scope.i = alt;
            }
            return true;
        },
        go: function(scope, __problems, callback) {
            problems.set(__problems);
            scope.count = __problems.length;

            var update = function(i) {
                scope.i = i;
                scope.problem = __problems[scope.i];
                moonboard.set(scope.problem.h);
                callback();

                if (__problems.length > perpage) {
                    scope.list = [];
                    var start = Math.min(scope.i, __problems.length - perpage);
                    scope.list = _.slice(__problems, start, start + perpage);
                } else {
                    scope.list = __problems;
                }
            };

            scope.ppage = function (event) {
                update(Math.max(scope.i - perpage, 0));
            };
            scope.prev = function (event) {
                if (scope.i > 0) {
                    update(scope.i-1);
                }
            };
            scope.rand = function (event) {
                update(Math.floor(Math.random() * __problems.length));
            };
            scope.next = function (event) {
                if (scope.i < (__problems.length - 1)) {
                    update(scope.i+1);
                }
            };
            scope.npage = function (event) {
                update(Math.min(scope.i + perpage, __problems.length - 1));
            };

            moonboard.load().then(
                function() {
                    update(Math.min(scope.i, __problems.length - 1));
                },
                function() {
                    scope.error = scope.error || { status: 500, data: 'Failed to load Moonboard' };
                }
            );
        }
    };
});