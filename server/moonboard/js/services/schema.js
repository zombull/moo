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
        'index.problems': '9e1d202bf94b6d43f07c54b276823091',
        'index.setters': 'a616c559e0b256277ddd70689dc5d4ac',
        problems: '06b620ffef389dcb6557dd8d8ad0a55c',
        setters: 'a4e271333ba9633b96185c3564d94ba3',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});