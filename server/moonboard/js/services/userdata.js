/**
 * Service for making queries to the database.
*/
moon.factory('userdata', function ($q, bug, storage, drive) {
    'use strict';

    var USER_DATA = 'zombull.moonboard.2016.userdata.json';

    var __drive = $q.defer();

    function initData() {
        return {
            projects: JSON.parse(localStorage.getItem('projects')) || {},
            ticks: JSON.parse(localStorage.getItem('ticks')) || {},
        };
    }

    // User data is stored on the user's Google Drive.  It's different from
    // the data provided by the server in that it is retrieved every time.
    // The major ramification of this is that there is a much higher chance
    // of not being able to connect, e.g. due to lack of network.  So, don't
    // completely freak out, we'll use the existing data in storage and will
    // forego updates to Google Drive (so we don't waste time).
    drive.get(USER_DATA, initData).then(
        function(file) {
            storage.set('projects', file.data.projects || {});
            storage.set('ticks', file.data.ticks || {});
            __drive.resolve(file.id);
        },
        function(error) {
            __drive.reject();
            console.log(error);
            alert("Unabled to retrieve user data from Google Drive");
        }
    );

    var ls = {
        get: function() {
            return JSON.parse(localStorage.getItem('pending'));
        },
        set: function(val) {
            localStorage.setItem('pending', JSON.stringify(val));
        }
    };

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
            __drive.promise.then(function(id) {
                var p = ls.get();
                ls.set({});
                if (!_.isEmpty(p)) {
                    drive.patch(id, p).then(
                        function() { },
                        function(error) {
                            bug.warn(true, error);
                            var existing = ls.get();
                            if (existing) {
                                _.each(p, function(patches, index) {
                                    _.each(patches, function(patch, key) {
                                        if (!existing[index].hasOwnProperty(key)) {
                                            existing[index][key] = patch;
                                        }
                                    });
                                });
                            } else {
                                existing = p;
                            }
                            ls.set(existing);
                        }
                    );
                }
            });
        }
    };
    pending.sync();

    return {
        get: function() {
            return __drive.promise;
        },
        add: function(index, data, key, val) {
            bug.on(!data.problems.hasOwnProperty(key));

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
