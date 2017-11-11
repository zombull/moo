host.factory('bug', function () {
    'use strict';

    return {
        bug: function(message) {
            console.trace();
            console.log(message);
            throw new Error("BUG: " + message);
        },
        on: function(condition, message) {
            if (condition) {
                console.trace();
                console.log(message);
                throw new Error("BUG: " + message);
            }
        },
        warn: function(condition, message) {
            if (condition) {
                console.trace();
                console.log(message);
                return true;
            }
            return false;
        }
    };
});