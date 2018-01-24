moon.factory('schema', function () {
    'use strict';

    var requests = {};

    var subdomains = { 'dark': {} };
    var metadata = {
        'index.problems': { subdomain: 'dark' },
        'index.setters': { local: true },
        images: { local: true },
        problems: { local: true },
        setters: { local: true },
        projects: { local: true, drive: true },
        ticks: { local: true, drive: true },
        exiles: { local: true, drive: true },
    };
    var checksums = {
        images: 'd12cd143a2d69484feaa72d1942bb979',
        'index.problems': 'ed37aecd0455f6e1540fe798761c314a',
        'index.setters': '8009c058736347e1c10336df35d0f1a3',
        problems: '2b689d6acdead9516313fd304b81c84c',
        setters: '0807df0a7942d1ddf40b9ed92bfc79a4',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});