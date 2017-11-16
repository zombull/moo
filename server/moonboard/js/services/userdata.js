/**
 * Service for making queries to the database.
*/
moon.factory('userdata', function ($q, bug, storage, drive) {
    'use strict';

    var USER_DATA = 'zombull.moonboard.2016.userdata.json';
    var USER_DATA_KEYS = ['projects', 'ticks', 'exiles'];

    var __drive = $q.defer();

    function initData() {
        var data = {};
        _.each(USER_DATA_KEYS, function(key) {
            data[key] = JSON.parse(localStorage.getItem(key)) || {};
        });
        return data;
    }

    var ls = {
        get: function() {
            return JSON.parse(localStorage.getItem('pending'));
        },
        set: function(val) {
            localStorage.setItem('pending', JSON.stringify(val));
        }
    };

    function syncPending(id, data) {
        var i = $q.defer();

        var p = ls.get();
        ls.set({});
        if (_.isEmpty(p)) {
            i.resolve(data);
        } else {
            drive.patch(id, p).then(
                function(patchedData) {
                    i.resolve(patchedData);
                },
                function(error) {
                    bug.warn(true, error);
                    var existing = ls.get();
                    if (existing) {
                        _.each(p, function(patches, index) {
                            _.each(patches, function(patch, key) {
                                if (!existing.hasOwnProperty(index) || !existing[index].hasOwnProperty(key)) {
                                    existing[index] = p[index] || {};
                                    existing[index][key] = patch;
                                }
                            });
                        });
                    } else {
                        existing = p;
                    }
                    ls.set(existing);
                    i.reject(error);
                }
            );
        }
        return i.promise;
    }

    function setPending(index, key, val) {
        var p = ls.get();
        p[index] = p[index] || {};
        p[index][key] = { val: val, add: !!val };
        ls.set(p);
    }

    var pending = {
        add: setPending,
        rm: setPending,
        sync: function() {
            __drive.promise.then(syncPending, function() { });
        }
    };

    // User data is stored on the user's Google Drive.  It's different from
    // the data provided by the server in that it is retrieved every time.
    // The major ramification of this is that there is a much higher chance
    // of not being able to connect, e.g. due to lack of network.  So, don't
    // completely freak out, we'll use the existing data in storage and will
    // forego updates to Google Drive (so we don't waste time).
    drive.get(USER_DATA, initData).then(
        function(file) {
            syncPending(file.id, file.data).then(
                function(data) {
                    _.each(USER_DATA_KEYS, function(key) {
                        storage.set(key, data[key] || {});
                    });
                    __drive.resolve(file.id);
                },
                function(error) {
                    __drive.reject();
                    alert("Unabled to sync user data to Google Drive");
                }
            );
        },
        function(error) {
            __drive.reject();
            alert("Unabled to retrieve user data from Google Drive");
        }
    );

    return {
        get: function() {
            return __drive.promise;
        },
        add: function(index, data, key, val) {
            data[index][key] = val;
            storage.set(index, data[index]);
            pending.add(index, key, val);
            return pending;
        },
        rm: function(index, data, key) {
            if (data[index].hasOwnProperty(key)) {
                delete data[index][key];
                storage.set(index, data[index]);
            }
            pending.rm(index, key);
            return pending;
        }
    };
});
