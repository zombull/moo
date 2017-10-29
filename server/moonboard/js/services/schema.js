moon.factory('schema', function () {
    'use strict';

    var requests = {};

    var subdomains = { 'dark': {}, 'side': {} };
    var metadata = {
        dark: { subdomain: 'dark' },
        side: { subdomain: 'side' },
        images: { local: true },
        problems: { local: true },
        setters: { local: true },
        ticks: { local: true },
        tocks: { local: true, ephemeral: true },
        projects: { local: true, ephemeral: true }
    };
    var checksums = {
        images: 'd12cd143a2d69484feaa72d1942bb979',
        index: '166c351ac262d5e2c19e144f840c4798',
        problems: 'add80bfd82cfe295bcaff49534af3584',
        setters: '484f12d3e02074807a00770584fa11a6',
        ticks: '4c9399e05de09d571ba02dda7158337c',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});