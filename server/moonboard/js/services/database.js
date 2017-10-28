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

    function getEverything(callback) {
        storage.get('master', getStorage.bind(this, function(data) {
            storage.get('ticks', getStorage.bind(this, function(ticks) {
                storage.get('tocks', getStorage.bind(this, function(tocks) {
                    storage.get('projects', getStorage.bind(this, function(projects) {
                        callback(data, ticks, tocks, projects);
                    }));
                }));
            }));
        }));
    }

    getEverything(function(data, ticks, tocks, projects){
        if (!data.hasOwnProperty('g') || update) {
            data.g = []
            for (var g = 0; g < 18; g++) {
                data.g[g] = [];
            }

            var end = _.size(data.p);
            _.each(data.i, function(problem, i) {
                if (i < end) {
                    problem.t = null;
                    if (ticks.hasOwnProperty(i)) {
                        problem.t = ticks[i];
                        if (tocks && tocks.hasOwnProperty(problem.u)) {
                            delete tocks[problem.u];
                        }
                    } else if (tocks && tocks.hasOwnProperty(problem.u)) {
                        problem.t = tocks[problem.u];
                    }
                    problem.g = problem.t ? problem.t.g : problem.g;
                    problem.s = (problem.t && problem.t.s) ? problem.t.s : problem.s;
                    problem.v = grades.convert(problem.g);

                    bug.on((problem.v/10) > 17, "Really, a V18?  Hello, Nalle!")
                    data.g[problem.v/10].push(i);
                }
            });
            if (tocks) {
                storage.set('tocks', tocks);
            }
            storage.set('master', data);
        }

        data.projects = projects || {};
        __data.resolve(data);
    });

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
                    bug.on(!data.p.hasOwnProperty(problem));
                    data.projects[problem] = project;
                    storage.set('projects', data.projects);
                });
            },
        },
        tock: {
            add: function(tock, scope, callback) {
                getData(scope, function(data) {
                    storage.get('tocks',function(tocks, error) {
                        if (error) {
                            scope.error = scope.error || error;
                        } else {
                            bug.on(!data.p.hasOwnProperty(tock.p));
                            var problem = data.i[data.p[tock.p]];
                            bug.on(problem.t !== null);
                            problem.t = tock;
                            problem.g = tock.g
                            problem.s = tock.s;
                            problem.v = grades.convert(problem.g);
                            storage.set('master', data);

                            bug.on(tocks && tocks.hasOwnProperty(tock.p));
                            tocks = tocks || {};
                            tocks[tock.p] = tock;
                            storage.set('tocks', tocks);
                            callback();
                        }
                    });
                });
            },
        }
    };
});