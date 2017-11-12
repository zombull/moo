moon.factory('history', function () {
    'use strict';

    var ls = {
        get: function() {
            return JSON.parse(localStorage.getItem('history'));
        },
        set: function(val) {
            localStorage.setItem('history', JSON.stringify(val));
        }
    };

    var h = ls.get() || {};
    return {
        get: function(p, alt) {
            return h.hasOwnProperty(p) ? h[p] : alt;
        },
        set: function(p, val) {
            h[p] = val;
            ls.set(h);
        },
    };
});