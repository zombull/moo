host.factory('grades', function () {
    'use strict';

    var conversions = {
        'VB':        0,
        'V0-':       1,
        'V0':        2,
        'V0+':       3,
        'V1':        10,
        'V2':        20,
        'V3':        30,
        'V4':        40,
        'V5':        50,
        'V6':        60,
        'V7':        70,
        'V8':        80,
        'V9':        90,
        'V10':       100,
        'V11':       110,
        'V12':       120,
        'V13':       130,
        'V14':       140,
        'V15':       150,
        'V16':       160,
        'V17':       170,
        'V18':       180,
        'V19':       190,
        'V20':       200,
        '3rd Class': 3000,
        '4th Class': 4000,
        '5.0':       5000,
        '5.1':       5010,
        '5.2':       5020,
        '5.3':       5030,
        '5.4':       5040,
        '5.5':       5050,
        '5.6':       5060,
        '5.7':       5070,
        '5.8':       5080,
        '5.8+':      5081,
        '5.9-':      5090,
        '5.9':       5091,
        '5.9+':      5092,
        '5.10a':     5100,
        '5.10b':     5101,
        '5.10c':     5102,
        '5.10d':     5103,
        '5.11a':     5110,
        '5.11b':     5111,
        '5.11c':     5112,
        '5.11d':     5113,
        '5.12a':     5120,
        '5.12b':     5121,
        '5.12c':     5122,
        '5.12d':     5123,
        '5.13a':     5130,
        '5.13b':     5131,
        '5.13c':     5132,
        '5.13d':     5133,
        '5.14a':     5140,
        '5.14b':     5141,
        '5.14c':     5142,
        '5.14d':     5143,
        '5.15a':     5150,
        '5.15b':     5151,
        '5.15c':     5152,
        '5.15d':     5153,
        '5.16a':     5160,
        '5.16b':     5161,
        '5.16c':     5162,
        '5.16d':     5163,
        '5.17a':     5170,
        '5.17b':     5171,
        '5.17c':     5172,
        '5.17d':     5173,
    };

    function convert(grade, min) {
        if (grade) {
            if (conversions.hasOwnProperty(grade)) {
                return conversions[grade];
            }
            if (/^5\.1\d$/.test(grade)) {
                if (min) {
                    return convert(grade + 'a');
                }
                else {
                    return convert(grade + 'd');
                }
            }
        }

        return undefined;
    }

    return {
        compare: function(min, max) {
            if ((!min || min[0] === 'v') && (!max || max[0] === 'v')) {
                min = min ? convert(min.toUpperCase(), true) : conversions.VB;
                max = max ? convert(max.toUpperCase(), false) : conversions.V20;
            }
            else if ((!min || min[0] === '5') && (!max || max[0] === '5')) {
                min = min ? convert(min, true) : conversions['3rd Class'];
                max = max ? convert(max, false) : conversions['5.16d'];
            }
            else {
                min = max = undefined;
            }
    
            if (min && max) {
                return function(grade) {
                    return grade >= min && grade <= max;
                };
            }
            return function() {
                return false;
            };
        },
        convert: convert
    };
});