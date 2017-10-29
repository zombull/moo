moon.factory('inspector', function ($location, $q, database, problems, calculator, grades, truthiness) {
    'use strict';

    var filter = function(options, index) {
        if (options.benchmark !== null || options.ticked !== null || options.grade || options.ascents ||
            options.stars || options.query || options.setby || options.setter) {

            options.query = options.query.replace(/^\s+/, '');
            return index.filter(function(entry) {
                /*
                 * Short circuit employed, keep less expensive operations early
                 * and move more expensive operations to the end, i.e. boolean
                 * checks first, string checks last.
                 */
                return  (options.benchmark === null || entry.hasOwnProperty('b') === options.benchmark) &&
                        (options.ticked === null || (entry.t !== null) === options.ticked) &&
                        (!options.grade || options.grade(entry.v)) &&
                        (!options.ascents || options.ascents(entry.a)) &&
                        (!options.stars || options.stars(entry.s)) &&
                        (!options.setby || options.setby.hasOwnProperty(entry.e)) &&
                        (!options.query || entry.l.indexOf(options.query) !== -1);
            });
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
        benchmark: /\s+(\!|@)b/,
        ticked: /\s+(\!|@)t/,
        setby: /\s+(\!|@)y\s?(\w+)/,
        setter: /\s+(@)e/,
        grade: /\s+(?:=|@)(v1\d|v\d)/,
        minGrade: /\s+>(v1\d|v\d)/,
        maxGrade: /\s+<(v1\d|v\d)/,
        ascents: /\s+(?:a=|@a)(\d+)/,
        minAscents: /\s+a>(\d+)/,
        maxAscents: /\s+a<(\d+)/,
        stars: /\s+(?:s=|@s)(\d+)/,
        minStars: /\s+s>(\d+)/,
        maxStars: /\s+s<(\d+)/,
    };

    function processRegEx(options, regEx) {
        var match = options.query.match(regEx);
        if (match) {
            options.query = options.query.replace(regEx, '');
            return match[1].toLowerCase() + (match[2] ? match[2].toLowerCase() : '');
        }
        return null;
    }

    return {
        search: function (query) {
            var deferred = $q.defer();
            var autoclear = false;
            if (query) {
                var min, max;
                var options = { query: ' ' + query.toLowerCase() };

                options.benchmark = truthiness(processRegEx(options, regExs.benchmark));
                options.ticked = truthiness(processRegEx(options, regExs.ticked));
                options.setby = truthiness(processRegEx(options, regExs.setby));
                options.setter = truthiness(processRegEx(options, regExs.setter));

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
                        deferred.resolve(filter(options, problems.get()));
                    });
                } else {
                    deferred.resolve(filter(options, problems.get()));
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