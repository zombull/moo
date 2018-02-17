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
        'index.problems': '2385e473fa2b22472bb4855dc3ea2a6e',
        'index.setters': '4ecf0dce3bd8aa74eaa2fd4c98586f2b',
        problems: '2ae41acd4764948a78329c4153aef0f4',
        setters: '43cf0e6b6750eaf00ca768a5731e2c24',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});