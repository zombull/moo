moon.factory('schema', function () {
    'use strict';

    var requests = {};

    var subdomains = {};
    var metadata = { master: { local: true }, ticks: { local: true }, projects: { local: true, ephemeral: true } };
    var checksums = { master: 'fa707ba5523acd49afaa244f7453fb38', ticks: '061b88ba97ac432d2bef39bf70b965c8' };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});