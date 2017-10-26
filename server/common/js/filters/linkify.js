host.filter('linkify', function($sce) {

    'use strict';

    return function(entry, external) {
        var target =  external ? ' target="_blank"' : '';
        var meta = entry.hasOwnProperty('g') ? '<i>&nbsp;&nbsp;({0})</i>'.format(entry.g) : '';
        var stars = ''
        if (entry.hasOwnProperty('g')) {
            stars = '&nbsp;&nbsp;'
            if (entry.hasOwnProperty('s')) {
                _.times(entry.s, function() {
                    stars += '&#x2605;';
                });
            } else {
                stars += '&#x2620;';
            }
        }
        var ticked = entry.hasOwnProperty('t') && entry.t ? '&#x2739;&nbsp;' : '';
        var benchmark = entry.hasOwnProperty('b') && entry.b ? '&#x272a;&nbsp;' : '';
        var html = '<a href="{0}"{1}>{2}{3}{4}{5}{6}</a>'.format(entry.u, target, ticked, benchmark, entry.n, meta, stars);
        return $sce.trustAsHtml(html);
    };
});
