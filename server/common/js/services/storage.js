host.factory('storage', function ($http, $q, bug, schema) {
    'use strict';

    // The iframe may not be listening for our first event, even if we wait until onload,
    // as being loaded does not mean it has run its initial script code.  Wait for a ping
    // from the subdomain, using a promise to track when the subdomain is ready.  Then()
    // callbacks are invoked in the order they are attached, so we can create postMessage
    // in this dedicated code, allowing pending schema.requests to assume postMessage is valid
    // once the promise is resolved.
    _.each(schema.subdomains, function(subdomain, name) {
        subdomain.ping = $q.defer();
        subdomain.pong = subdomain.ping.promise;
        subdomain.pong.then(
            function(source) {
                subdomain.window = document.getElementById(name).contentWindow;
                subdomain.postMessage = function(message) {
                    this.window.postMessage(message, 'http://' + name + '.zombull.xyz:3000');
                };
                subdomain.postMessage('pong');
            },
            function(error) {
                // Completely hosed if the ping is rejected.
                bug.bug(error);
            }
        );
    });

    function onMessage(event) {
        var match = event.origin.match(/^http:\/\/([A-Za-z0-9-]+)\.zombull\.xyz:3000$/);
        if (match && schema.subdomains.hasOwnProperty(match[1])) {
            if (event.data === 'ping') {
                if (!schema.subdomains[match[1]].postMessage) {
                    schema.subdomains[match[1]].ping.resolve();
                }
            }
            else {
                if (schema.requests.hasOwnProperty(event.data.key)) {
                    _.each(schema.requests[event.data.key], function(request) {
                        if (event.data.value) {
                            request.resolve(event.data.value);
                        }
                        else {
                            request.reject();
                        }
                    });
                    delete schema.requests[event.data.key];
                }
            }
        }
    }
    window.addEventListener('message', onMessage, false);

    function postMessage(subdomain, message, request) {
        if (schema.subdomains[subdomain].postMessage) {
            schema.subdomains[subdomain].postMessage(message);
        }
        else {
            schema.subdomains[subdomain].pong.then(
                function(source) {
                    schema.subdomains[subdomain].postMessage(message);
                },
                function(error) {
                    // Completely hosed if the ping is somehow rejected.
                    bug.bug(error);
                }
            );
        }
    }

    var substorage = {
        get: function(key) {
            var request = $q.defer();
            if (!schema.metadata.hasOwnProperty(key)) {
                request.reject(true);
            }
            else {
                if (schema.metadata[key].local) {
                    var local = localStorage.getItem(key);
                    if (local) {
                        request.resolve(local);
                    }
                    else {
                        request.reject();
                    }
                }
                else {
                    if (schema.requests.hasOwnProperty(key)) {
                        schema.requests[key].push(request);
                    }
                    else {
                        schema.requests[key] = [];
                        schema.requests[key].push(request);

                        postMessage(schema.metadata[key].subdomain, { method: 'get', key: key });

                        // Don't really want to fall back to the server, this code needs to be rock solid.
                        // $timeout(request.reject());
                    }
                }
            }
            return request.promise;
        },

        set: function(key, value) {
            if (schema.metadata.hasOwnProperty(key)) {
                if (schema.metadata[key].local) {
                    localStorage.setItem(key, value);
                }
                else {
                    postMessage(schema.metadata[key].subdomain, { method: 'set', key: key, value: value });
                }
            }
        }
    };

    var cache = {
        data: {},
        checksums: JSON.parse(localStorage.getItem('checksums')) || { },
    };

    function deposit(name, value, update) {
        // Do not overwrite existing data unless explicit requested to do so
        // as part of an update.  This prevents overwriting an update with
        // stale data from local storage.
        if (update || !cache.data.hasOwnProperty(name)) {
            cache.data[name] = value;
            if (update) {
                  substorage.set(name, JSON.stringify(value));
            }
        }

        if (update || !cache.checksums.hasOwnProperty(name)) {
            cache.checksums[name] = schema.checksums[name];
            localStorage.setItem('checksums', JSON.stringify(cache.checksums));
        }
    }

    function doUpdate(name, callback) {
        // Append the hash as a query string to create a unique URI.  This allows
        // CDNs to cache the data but guarantees we'll get the latest version, all
        // without having to store multiple versions on the server.
        $http.get('data/' + name + '?version=' + cache.checksums[name]).then(
            function(response) {
                deposit(name, response.data, true);
                callback(response.data);
            },
            function(response) {
                callback(null, response);
            }
        );
    }

    return {
        checksums: function() {
            return cache.checksums;
        },
        get: function(name, callback) {
            if (cache.data.hasOwnProperty(name)) {
                callback(cache.data[name]);
            }
            else {
                substorage.get(name).then(
                    function(value) {
                        var d = JSON.parse(value);
                        deposit(name, d, false);
                        callback(d);
                    },
                    function(error) {
                        if (error) {
                            callback(null, error);
                        } else {
                            doUpdate(name, callback);
                        }
                    }
                );
            }
        },
        set: function(name, data) {
            deposit(name, data, true);
        },
        update: function(name) {
            doUpdate(name, function() { });
        }
    };
});