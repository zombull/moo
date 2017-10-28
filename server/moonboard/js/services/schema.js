moon.factory('schema', function () {
    'use strict';

    var requests = {};

    var subdomains = {};
    var metadata = {
        master: { local: true },
        ticks: { local: true },
        tocks: { local: true, ephemeral: true },
        projects: { local: true, ephemeral: true }
    };
    var checksums = { master: 'fa707ba5523acd49afaa244f7453fb38', ticks: '4c9399e05de09d571ba02dda7158337c' };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});