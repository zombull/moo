moon.factory('schema', function () {
    'use strict';

    var requests = {};

    var subdomains = { 'darkZZZZ': {} };
    var metadata = {
        'index.problems': { subdomain: 'darkZZZZ' },
        'index.setters': { local: true },
        images: { local: true },
        problems: { local: true },
        setters: { local: true },
        projects: { local: true, drive: true },
        ticks: { local: true, drive: true },
        exiles: { local: true, drive: true },
    };

    var checksums = SCHEMA_CHECKSUMS;

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});