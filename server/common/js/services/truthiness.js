host.factory('truthiness', function () {
    'use strict';

    return function(b, v) {
        if (!v) {
            return b === '!' ? false : true;
        }
        return {
            b: b !== '!',
            v: v.trim()
        };
    };
});