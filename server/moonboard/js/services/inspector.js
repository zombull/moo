moon.factory('inspector', function ($location, $q, bug, database, problems, calculator, grades, truthiness) {
    'use strict';

    // Keep track of the last results, this is used when browsing
    // search results to remember the last search, especially when
    // bouncing around the history.
    var __results = [];

    var filter = function(options, index, browsing) {
        if (options.benchmark !== null || options.ticked !== null || options.project !== null ||
            options.todo !== null || options.exiled !== null || options.upgraded != null ||
            options.downgraded !== null ||
            options.grade || options.ascents || options.stars || options.query || options.setby ||
            options.setter || options.holds || options.noholds) {

            options.query = options.query.replace(/^\s+/, '');
            __results = index.filter(function(entry) {
                /*
                 * Short circuit employed, keep less expensive operations early
                 * and move more expensive operations to the end, i.e. boolean
                 * checks first, string checks last.
                 */
                return  (options.benchmark === null || entry.hasOwnProperty('b') === options.benchmark) &&
                        (options.ticked === null || (entry.t !== null) === options.ticked) &&
                        (options.project === null || (entry.p !== null && (entry.p.s >= 2 || entry.p.a >= 4) === options.project)) &&
                        (options.todo === null || (entry.p !== null && (entry.p.s < 2 && entry.p.a < 4) === options.todo)) &&
                        (options.exiled === null || entry.e === options.exiled) &&
                        (options.upgraded === null || (entry.hasOwnProperty('w') && entry.w < entry.v === options.upgraded)) &&
                        (options.downgraded === null || (entry.hasOwnProperty('w') && entry.w > entry.v === options.downgraded)) &&
                        (!options.grade || options.grade(entry.v)) &&
                        (!options.ascents || options.ascents(entry.a)) &&
                        (!options.stars || options.stars(entry.s)) &&
                        (!options.setby || options.setby.hasOwnProperty(entry.r)) &&
                        (!options.query || entry.l.indexOf(options.query) !== -1) &&
                        (!options.holds || options.holds(entry.h)) &&
                        (!options.noholds || options.noholds(entry.h));
            });
            if (!browsing && __results.length > 1) {
                __results.unshift({ n: '*** Browse Results ***', u: 'q/' + encodeURI(options.raw) });
            }
            return __results;
        }
        return index;
    };

    var filterSetters = function(options, index) {
        if (options.query) {
            options.query = options.query.replace(/^\s+/, '');
            return index.filter(function(entry) {
                return (entry.l.indexOf(options.query) !== -1);
            });
        }
        return index;
    };

    var regExs = {
        benchmark: /\s+(\!|\.)b/,
        ticked: /\s+(\!|\.)t/,
        project: /\s+(\!|\.)p/,
        todo: /\s+(\!|\.)o/,
        exiled: /\s+(\!|\.)e/,
        upgraded: /\s+(\!|\.)u/,
        downgraded: /\s+(\!|\.)d/,
        setby: /\s+(\!|\.)r\s?(\w+)/,
        setter: /\s+(\.)n/,
        grade: /\s+(?:=|\.)(v1\d|v\d)/,
        holds: /\s+(?:h=|\.h)((?:[a-kA-K][0-9][0-9]?)(?:,[a-kA-K][0-9][0-9]?)*)/,
        noholds: /\s+(?:h!=|!h)((?:[a-kA-K][0-9][0-9]?)(?:,[a-kA-K][0-9][0-9]?)*)/,
        minGrade: /\s+(v1\d|v\d)\.\./,
        maxGrade: /\s+\.\.(v1\d|v\d)/,
        ascents: /\s+(?:a=|\.a)(\d+)/,
        minAscents: /\s+a(\d+)\.\./,
        maxAscents: /\s+\.\.a(\d+)/,
        stars: /\s+(?:s=|\.s)(\d+)/,
        minStars: /\s+s(\d+)\.\./,
        maxStars: /\s+\.\.s(\d+)/,
    };

    function processRegEx(options, regEx, fn) {
        var match = options.query.match(regEx);
        if (match) {
            options.query = options.query.replace(regEx, '');
            if (fn) {
                return fn(match[1], match[2]);
            }
            return match[1].toLowerCase();
        }
        return null;
    }

    function filterHolds(b, q) {
        q = q.toUpperCase().split(',');
        return function(holds) {
            for (var i = 0; i < q.length; i++) {
                if ((holds.indexOf(q[i]) !== -1) !== b) {
                    return false;
                }
            }
            return true;
        };
    }

    return {
        search: function (query, browsing) {
            if (browsing) {
                if (__results.length > 0 && !__results[0].hasOwnProperty('l')) {
                    __results.shift();
                }
                problems.set(__results);
            }

            var deferred = $q.defer();
            if (query) {
                var min, max;
                var options = { raw: query.toLowerCase(), query: ' ' + query.toLowerCase() };

                options.benchmark = processRegEx(options, regExs.benchmark, truthiness);
                options.ticked = processRegEx(options, regExs.ticked, truthiness);
                options.project = processRegEx(options, regExs.project, truthiness);
                options.todo = processRegEx(options, regExs.todo, truthiness);
                options.exiled = processRegEx(options, regExs.exiled, truthiness);
                options.upgraded = processRegEx(options, regExs.upgraded, truthiness);
                options.downgraded = processRegEx(options, regExs.downgraded, truthiness);
                options.setby = processRegEx(options, regExs.setby, truthiness);
                options.setter = processRegEx(options, regExs.setter, truthiness);

                options.holds = processRegEx(options, regExs.holds, filterHolds.bind(null, true));
                options.noholds = processRegEx(options, regExs.noholds, filterHolds.bind(null, false));

                min = max = processRegEx(options, regExs.grade);
                if (!min) {
                    min = processRegEx(options, regExs.minGrade);
                    max  = processRegEx(options, regExs.maxGrade);
                }
                if (min || max) {
                    options.grade = grades.compare(min, max);
                }

                min = max = processRegEx(options, regExs.ascents);
                if (!min) {
                    min = processRegEx(options, regExs.minAscents);
                    max = processRegEx(options, regExs.maxAscents);
                }
                if (min || max) {
                    options.ascents = calculator(min, max);
                }

                min = max = processRegEx(options, regExs.stars);
                if (!min) {
                    min = processRegEx(options, regExs.minStars);
                    max = processRegEx(options, regExs.maxStars);
                }
                if (min || max) {
                    options.stars = calculator(min, max);
                }

                if (options.setter) {
                    database.setters(function(setters) {
                        deferred.resolve(filterSetters(options, setters));
                    });
                } else if (options.setby) {
                    database.setterIds(function(setters) {
                        // Build an object whose properties are setter ids.
                        // Filter will check a problem's setter ID against
                        // the object's properties.  @key below is the URL
                        // of the setter, i.e. /s/<name> minus spaces and
                        // lowercased.
                        var setby = options.setby;
                        options.setby = {};
                        _.each(setters, function(id, key) {
                            if ((key.indexOf(setby.v) !== -1) === setby.b) {
                                options.setby[id] = true;
                            }
                        });
                        deferred.resolve(filter(options, problems.get(), browsing));
                    });
                } else {
                    deferred.resolve(filter(options, problems.get(), browsing));
                }
            }
            else {
                deferred.resolve(problems.get());
            }
            return deferred.promise;
        },
        autoclear: function(query) {
            var clear = true;
            var options = { query: ' ' + query.toLowerCase() };
            _.each(regExs, function(regex) {
                    if (clear && processRegEx(options, regex)) {
                        clear = false;
                    }
            });
            return clear;
        }
    };
});