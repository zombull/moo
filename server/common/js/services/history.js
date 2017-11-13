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
    var e = {};
    return {
        get: function(p, alt) {
            return h.hasOwnProperty(p) ? h[p] : e.hasOwnProperty(p) ? e[p] : alt;
        },
        set: function(p, val, ephemeral) {
            if (!ephemeral) {
                h[p] = val;
                ls.set(h);
            } else {
                e[p] = val;
            }
        },
    };
});