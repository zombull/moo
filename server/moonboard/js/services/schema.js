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
        images: 'a689dbfef34cdd7bbc2fd6a27f6d5075',
        'index.problems': '09f9a1767c2a3b2764bdc7c0591a31f4',
        'index.setters': 'afdf0006a9f4dbc479facbf587af94dc',
        problems: 'c34adcad3665627a3942b3f0969aece4',
        setters: 'fdf1d9d495c8ad90fd35432f8ea1b92c',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});