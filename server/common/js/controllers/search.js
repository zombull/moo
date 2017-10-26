/**
 *
 */
host.controller('SearchController', function SearchController($location, $timeout, inspector) {
    'use strict';

    var self = this;

    var FORWARD_SLASH = 191;
    var focusSearch = function(event) {
        // Make sure the user is not typing into an input.  No modifier is allowed.
        if (event.keyCode === FORWARD_SLASH && !event.shiftKey && !event.ctrlKey && !event.altKey && !event.metaKey && event.target.tagName.toLowerCase() !== 'input') {
            document.getElementById('search').focus();

            event.preventDefault();
            event.stopPropagation();
        }
    };
    window.addEventListener('keydown', focusSearch, false);

    self.query = '';
    self.search = inspector.search;

    self.searchTextChange = function(query) {
        self.query = query;
    }

    function clearSearch() {
        if (inspector.autoclear(self.query)) {
            self.searchText = '';
        } else {
            self.searchText = self.query;
        }
    }

    var onFocus;
    onFocus = function(event) {
        // Conditionally clear the search text when the input is focused.
        // This needs to be done in $timeout, adjusting searchText while
        // Angular is doing its thing will make it think the selection
        // changed and Angular will display the list again and again and again.
        $timeout(clearSearch);

        // Remove the event handler, we want to preserve the user's input if the focus on something
        // else prior to selecting an element.
        document.getElementById('search').removeEventListener('focus', onFocus, false);
    };

    self.selectedItemChange = function(item) {
        if (item && item.u) {
            $location.path(item.u);

            // Focus on the main column to hide any soft keyboard.
            document.getElementById('main').focus();

            // Add the event handler to clear the input the next time it is focused.  Do this in
            // a timeout so that we don't try clearing the input until Angular has processed the
            // selection.
            $timeout(function() {
                document.getElementById('search').addEventListener('focus', onFocus, false);
            });
        }
    };
});