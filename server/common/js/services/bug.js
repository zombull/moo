host.factory('bug', function () {
    'use strict';

    return {
        bug: function(message) {
            console.trace();
            console.log(message);
        },
        on: function(condition, message) {
            if (condition) {
                console.trace();
                console.log(message);
            }
        }
    };
});