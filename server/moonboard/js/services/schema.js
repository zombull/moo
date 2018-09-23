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
        'index.problems': 'e3d7fad56fda1286c14e4d5ef7417a7f',
        'index.setters': '46042e647222019c096a4f03670d8564',
        problems: '08fb85dad8d2fad46446633b6496bb0b',
        setters: 'd1f5db99807f68e909ac82028226c70d',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});