/**
 * Service for making queries to the database.
*/
moon.factory('database', function ($http, $q, bug, grades, storage, userdata, schema) {
    'use strict';

    var __data = $q.defer();
    var __updates = [];

    var checksums = storage.checksums();
    _.each(schema.checksums, function(sum, key) {
        if (checksums[key] !== sum) {
            // Append the hash as a query string to create a unique URI.  This
            // allows the data to be cached, e.g. by proxies, but guarantees
            // we'll get the latest version,  And we don't have to store multiple
            // versions on the server.
            var i = $q.defer();
            $http.get('data/' + key + '?version=' + sum).then(
                function(response) {
                    storage.set(key, response.data);
                    i.resolve();
                },
                function(error) {
                    i.reject(error);
                }
            );
            __updates.push(i.promise);
        }
    });

    function unpack(data) {
        bug.on(data.hasOwnProperty('grades'));
        bug.on(data.hasOwnProperty('index'));

        data.index = {
            problems: data['index.problems'],
            setters: data['index.setters'],
        };
        delete data['index.problems'];
        delete data['index.setters'];

        data.grades = [];
        for (var g = 0; g < 18; g++) {
            data.grades[g] = [];
        }

        _.each(data.index.problems, function(problem, i) {
            problem.t = null;
            if (data.ticks.hasOwnProperty(problem.u)) {
                problem.t = data.ticks[problem.u];
                if (problem.t.g != problem.g) {
                    // Save the original grade if it is different than the ticked grade.
                    problem.o = problem.g;
                    problem.w = grades.convert(problem.g);
                }
                problem.g = problem.t.g;
                problem.s = problem.t.s ? problem.t.s : problem.s;

                // User data only modifies the specified index, i.e.
                // this is safe even though we're still unpacking.
                if (bug.warn(data.projects.hasOwnProperty(problem.u), "Project exists for tick, queueing removal")) {
                    userdata.rm('projects', data, problem.u).sync();
                }
                if (bug.warn(data.exiles.hasOwnProperty(problem.u), "Exile exists for tick, queueing removal")) {
                    userdata.rm('exiles', data, problem.u).sync();
                }
            }
            problem.p = data.projects.hasOwnProperty(problem.u) ? data.projects[problem.u] : null;
            problem.e = data.exiles.hasOwnProperty(problem.u);
            problem.v = grades.convert(problem.g);

            bug.on((problem.v/10) > 17, "Really, a V18?  Hello, Nalle!");
            data.grades[problem.v/10].push(i);
        });
        __data.resolve(data);
    }

    function getStorage(callback, data, error) {
        if (error) {
            __data.reject(error);
        } else {
            callback(data);
        }
    }

    function fetch() {
        var data = {};
        var keys = _.keys(schema.metadata);
        function __fetch(key) {
            if (key) {
                storage.get(key, getStorage.bind(null, function(val) {
                    data[key] = val;
                    __fetch(keys.shift());
                }));
            } else {
                unpack(data);
            }
        }
        __fetch(keys.shift());
    }

    $q.all(__updates).then(
        function() {
            // Syncing user data from Google Drive is not mandatory,
            // i.e. we can continue on even if syncing fails, we'll
            // just go into an 'offline' mode.
            userdata.get().then(fetch, fetch);
        },
        function(error) {
            __data.reject(error);
        }
    );

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

    function assertUserData(data, problemUrl, exist, notExist) {
        _.each(exist, function(index) {
            bug.on(!data[index].hasOwnProperty(problemUrl), problemUrl);
        });
        _.each(notExist, function(index) {
            bug.on(data[index].hasOwnProperty(problemUrl), problemUrl);
        });
    }

    function getProblem(data, problemUrl) {
        return data.index.problems[data.problems[problemUrl]];
    }

    return {
        all: function(callback, scope) {
            getData(scope, callback);
        },
        images: function(callback, scope) {
            getData(scope, function(data) {
                callback(data.images);
            });
        },
        setters: function(callback, scope) {
            getData(scope, function(data) {
                callback(data.index.setters);
            });
        },
        setterIds: function(callback, scope) {
            getData(scope, function(data) {
                callback(data.setters);
            });
        },
        project: {
            add: function(problemUrl, project, scope) {
                getData(scope, function(data) {
                    assertUserData(data, problemUrl, ['problems'], ['exiles', 'ticks']);
                    getProblem(data, problemUrl).p = project;
                    userdata.add('projects', data, problemUrl, project).sync();
                });
            },
            rm: function(problemUrl, scope) {
                getData(scope, function(data) {
                    assertUserData(data, problemUrl, ['problems'], []);
                    getProblem(data, problemUrl).p = null;
                    userdata.rm('projects', data, problemUrl).sync();
                });
            },
            migrate: function(problemUrl, scope) {
                getData(scope, function(data) {
                    assertUserData(data, problemUrl, ['problems', 'projects'], []);

                    var prob = getProblem(data, problemUrl);
                    var mig = data.projects[problemUrl];
                    userdata.rm('projects', data, problemUrl).sync();
                    userdata.add('projects', data, prob.moon, mig).sync();
                });
            }
        },
        tick: {
            add: function(problemUrl, tick, scope, callback) {
                getData(scope, function(data) {
                    assertUserData(data, problemUrl, ['problems'], ['exiles', 'ticks']);

                    var problem = getProblem(data, problemUrl);
                    problem.t = tick;
                    problem.g = tick.g;
                    problem.s = tick.s;
                    problem.v = grades.convert(problem.g);

                    userdata.rm('projects', data, problemUrl);
                    userdata.add('ticks', data, problemUrl, tick).sync();
                    callback();
                });
            },
            rm: function(problemUrl, scope) {
                getData(scope, function(data) {
                    assertUserData(data, problemUrl, ['problems', 'ticks'], []);
                    getProblem(data, problemUrl).t = null;
                    userdata.rm('ticks', data, problemUrl).sync();
                });
            },
            migrate: function(problemUrl, scope) {
                getData(scope, function(data) {
                    assertUserData(data, problemUrl, ['problems', 'ticks'], []);

                    var prob = getProblem(data, problemUrl);
                    var mig = data.ticks[problemUrl];
                    userdata.rm('ticks', data, problemUrl).sync();
                    userdata.add('ticks', data, prob.moon, mig).sync();
                });
            }
        },
        exile: {
            add: function(problemUrl, scope) {
                getData(scope, function(data) {
                    assertUserData(data, problemUrl, ['problems'], ['exiles', 'ticks']);
                    getProblem(data, problemUrl).e = true;
                    userdata.add('exiles', data, problemUrl, true).sync();
                });
            },
            rm: function(problemUrl, scope) {
                getData(scope, function(data) {
                    assertUserData(data, problemUrl, ['problems', 'exiles'], ['ticks']);
                    getProblem(data, problemUrl).e = false;
                    userdata.rm('exiles', data, problemUrl).sync();
                });
            },
            migrate: function(problemUrl, scope) {
                getData(scope, function(data) {
                    assertUserData(data, problemUrl, ['problems', 'exiles'], []);

                    var prob = getProblem(data, problemUrl);
                    var mig = data.exiles[problemUrl];
                    userdata.rm('exiles', data, problemUrl).sync();
                    userdata.add('exiles', data, prob.moon, mig).sync();
                });
            }
        }
    };
});