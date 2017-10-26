moon.factory('problems', function () {
    'use strict';

    var problems = [];

    return {
        get: function() {
            return problems;
        },
        set: function(p) {
            problems = p;
        },
        reset: function() {
            problems = [];
        }
    };
});