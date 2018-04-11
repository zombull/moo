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
        'index.problems': 'af2b0cc1bc51e6ae7db928e884d9b8e3',
        'index.setters': 'c7ea821c2156b4ebfe6c5c913a2f9d13',
        problems: '68438529e3c391cdc4ff865c5f7248dc',
        setters: 'ae31db31237c84624ba0e8ab7ee9f6b0',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});