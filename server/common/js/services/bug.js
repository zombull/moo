host.factory('bug', function () {
    'use strict';

    return {
        bug: function(message) {
            throw message;
        },
        on: function(condition, message) {
            if (condition) {
                throw message;
            }
        }
    };
});