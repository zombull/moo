/**
 * Service for making queries to the database.
*/
moon.factory('database', function ($q, bug, grades, storage, schema) {
    'use strict';

    var update = false;
    var __data = $q.defer();

    _.each(storage.checksums(), function(sum, key) {
        if (sum !== schema.checksums[key]) {
            storage.update(key);
            update = true;
        } else {
            storage.get(key, function() { });
        }
    });


    function getStorage(callback, data, error) {
        if (error) {
            __data.reject(error);
        } else {
            callback(data);
        }
    }

    storage.get('master', getStorage.bind(this, function(data) {
        storage.get('ticks', getStorage.bind(this, function(ticks) {
            if (!data.hasOwnProperty('g') || update) {
                data.g = []
                for (var g = 0; g < 18; g++) {
                    data.g[g] = [];
                }                
                // Unpack difficulty and tick info into problems, then update
                // storage.
                var end = _.size(data.p);
                _.each(data.i, function(problem, i) {
                    if (i < end) {
                        problem.t = ticks.hasOwnProperty(i) ? ticks[i] : null;
                        problem.g = problem.t ? problem.t.g : problem.g;
                        problem.s = (problem.t && problem.t.s) ? problem.t.s : problem.s;
                        problem.v = grades.convert(problem.g);

                        bug.on((problem.v/10) > 17, "Really, a V18?  Hello, Nalle!")
                        data.g[problem.v/10].push(i);
                    }
                });
            }
            storage.set('master', data);
            __data.resolve(data);
        }));
    }));

    var getData = function(scope, callback) {
        __data.promise.then(
            function(data) {
                callback(data);
            },
            function(error) {
                scope.error = scope.error || error;    
            }
        );
    };

    return {
        all: function(callback, scope) {
            getData(scope, callback);

        },
        images: function(callback, scope) {
            getData(scope, function(data) {
                callback(data.img);
            });
        },
        setters: function(callback, scope) {
            getData(scope, function(data) {
                callback(_.slice(data.i, _.size(data.p)));
            });
        },
        setterIds: function(callback, scope) {
            getData(scope, function(data) {
                callback(data.s);
            });
        }
    };
});