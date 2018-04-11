moon.factory('schema', function () {
    'use strict';

    var requests = {};

    var subdomains = { 'dark2016': {} };
    var metadata = {
        'index.problems': { subdomain: 'dark2016' },
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
        'index.problems': '76cdfacd568d75d85732b37d50631292',
        'index.setters': '228cad5216c84ce6f7ae23a9a5aabe3f',
        problems: '0a091b33b170f9a235b387a4a6fcb08b',
        setters: '06d44c489d80ed6bf62c523a64246755',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});