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
        'index.problems': 'ccaf7a296e5f6a393b2b82b2160b6f5a',
        'index.setters': '5ece79b95f8b9c6ece7a52a360fda4a6',
        problems: '3120a27771984a05c59352b66e7918c1',
        setters: 'dac5cdf0c3bf3ec47794702b486f103e',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});