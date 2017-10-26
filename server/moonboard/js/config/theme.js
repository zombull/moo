/**
 *
 */
moon.config(function($mdThemingProvider) {
    'use strict';

    var lightBlueZ = $mdThemingProvider.extendPalette('light-blue', { '500': '#00EFDE' });
    var deepOrangeZ = $mdThemingProvider.extendPalette('deep-orange', { 'A200': '#FF3900' });

    $mdThemingProvider.definePalette('lightBlue', lightBlueZ);
    $mdThemingProvider.definePalette('deepOrange', deepOrangeZ);

    $mdThemingProvider.theme('default')
        .primaryPalette('lightBlue')
        .accentPalette('deepOrange')
});