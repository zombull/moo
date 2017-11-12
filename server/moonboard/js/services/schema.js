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
    };
    var checksums = {
        images: 'd12cd143a2d69484feaa72d1942bb979',
        'index.problems': '1be0c5dc7ca570480c7e16635641cf83',
        'index.setters': 'f7a3e016aecb0794c60009d3178a7e86',
        problems: 'add80bfd82cfe295bcaff49534af3584',
        setters: '9517ff32978a2f84af205ed5e58f2ee6',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});