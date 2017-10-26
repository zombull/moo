/*global angular */
/*jshint unused:false */

// Define a prototype to provide C# style functionality for formatting a string.  To use, invoke on a
// string,  e.g. 'Your mother was a {0} and your father smelt of {1}'.format('hamster', 'elderberries').
if (!String.prototype.format) {
    String.prototype.format = function() {
        var args = arguments;
        return this.replace(/{(\d+)}/g, function(match, number) {
            return typeof args[number] != 'undefined' ? args[number] : match;
        });
    };
}

/**
 * Instantiate the Moon module, this global will be used throughout.
 *
 * @type {angular.Module}
 */
var moon = angular.module('moon', ['ngRoute', 'ngMaterial']);
var host = moon;