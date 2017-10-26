host.factory('truthiness', function () {
    'use strict';

    return function(val) {
        if (val === null) {
            return val
        }
        if (val.length === 1) {
            return val === '!' ? false : true;
        }
        return {
            b: val.substring(0, 1) !== '!',
            v: val.substring(1).trim()
        };
    };
});