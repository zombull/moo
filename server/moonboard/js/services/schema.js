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
        'index.problems': 'ef3ce02cdd3f6c4f6b0866847d97b121',
        'index.setters': 'dfb86f33530aa9eabdbb018f852cd825',
        problems: '6e1094b38af109cbcf7d1966d675c52c',
        setters: '2f8fabeb56d069b2117dba381093116a',
    };

    return {
        requests: requests,
        metadata: metadata,
        subdomains: subdomains,
        checksums: checksums,
    };
});