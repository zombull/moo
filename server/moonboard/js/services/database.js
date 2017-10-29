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
        }
    });

    function getStorage(callback, data, error) {
        if (error) {
            __data.reject(error);
        } else {
            callback(data);
        }
    }
    var data = {}
    var keys = _.keys(schema.metadata);
    function fetch(key) {
        if (key) {
            storage.get(key, getStorage.bind(this, function(val) {
                if (key === 'dark' || key === 'side') {
                    console.log(key + " = "+ (((JSON.stringify(val).length + 'index'.length)* 2)/1024).toFixed(2)+" KB");
                }
                data[key] = val;
                fetch(keys.shift());
            }));
        } else {
            bug.on(data.hasOwnProperty('grades'));

            data['index'] = data.dark.concat(data.side);
            delete data[dark];
            delete data[side];

            data.grades = []
            for (var g = 0; g < 18; g++) {
                data.grades[g] = [];
            }

            var end = _.size(data.problems);
            _.each(data.index, function(problem, i) {
                if (i < end) {
                    problem.t = null;
                    if (data.ticks.hasOwnProperty(i)) {
                        problem.t = data.ticks[i];
                        if (data.tocks && data.tocks.hasOwnProperty(problem.u)) {
                            delete data.tocks[problem.u];
                        }
                    } else if (data.tocks && data.tocks.hasOwnProperty(problem.u)) {
                        problem.t = data.tocks[problem.u];
                    }
                    if (problem.t && data.projects && data.projects.hasOwnProperty(problem.u)) {
                        delete data.projects[problem.u];
                    }
                    problem.g = problem.t ? problem.t.g : problem.g;
                    problem.s = (problem.t && problem.t.s) ? problem.t.s : problem.s;
                    problem.v = grades.convert(problem.g);

                    bug.on((problem.v/10) > 17, "Really, a V18?  Hello, Nalle!")
                    data.grades[problem.v/10].push(i);
                }
            });
            if (data.tocks) {
                storage.set('tocks', data.tocks);
            }
            if (data.projects) {
                storage.set('projects', data.projects);
            }

            data.projects = data.projects || {};
            data.tocks = data.tocks || {};
            __data.resolve(data);
        }
    }
    fetch(keys.shift());

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
                callback(data.images);
            });
        },
        setters: function(callback, scope) {
            getData(scope, function(data) {
                callback(_.slice(data.index, _.size(data.problems)));
            });
        },
        setterIds: function(callback, scope) {
            getData(scope, function(data) {
                callback(data.setters);
            });
        },
        project: {
            get: function(problem, scope, callback) {
                getData(scope, function(data) {
                    if (data.projects.hasOwnProperty(problem)) {
                        callback(data.projects[problem]);
                    } else {
                        callback(null);
                    }
                });
            },
            set: function(problem, project, scope) {
                getData(scope, function(data) {
                    bug.on(!data.problems.hasOwnProperty(problem));
                    data.projects[problem] = project;
                    storage.set('projects', data.projects);
                });
            },
            rm: function(problem, scope) {
                getData(scope, function(data) {
                    if (data.projects.hasOwnProperty(problem)) {
                        delete data.projects[problem];
                        storage.set('projects', data.projects);
                    }
                });
            }
        },
        tock: {
            add: function(tock, scope, callback) {
                getData(scope, function(data) {
                    bug.on(!data.problems.hasOwnProperty(tock.p));
                    var problem = data.index[data.problems[tock.p]];
                    bug.on(problem.t !== null);

                    if (data.projects.hasOwnProperty(tock.p)) {
                        delete data.projects[tock.p];
                        storage.set('projects', data.projects);
                    }

                    problem.t = tock;
                    problem.g = tock.g
                    problem.s = tock.s;
                    problem.v = grades.convert(problem.g);

                    bug.on(data.tocks.hasOwnProperty(tock.p));
                    data.tocks[tock.p] = tock;
                    storage.set('tocks', data.tocks);
                    callback();
                });
            },
        }
    };
});