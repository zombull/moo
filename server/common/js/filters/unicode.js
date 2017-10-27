host.filter('unicode', function($sce) {
    'use strict';

    return function(input){
        return $sce.trustAsHtml('&#x{0};'.format(input));
    }
});
    