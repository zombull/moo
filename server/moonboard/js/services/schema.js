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
        'index.problems': '5008bb7dfb37e732c95f647b3bad83ef',
        'index.setters': 'f8b161f04dbbd58293e6068b103328bc',
        problems: '3cfc4c0b9cf30d23c37145ed1cffd55d',
        setters: '2980da4a154fef74f65558b21060a153',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});