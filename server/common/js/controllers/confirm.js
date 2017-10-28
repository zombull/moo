moon.controller('ConfirmController', function($scope, $mdDialog, $timeout, prompt, buttons) {
    $scope.prompt = prompt;
    $scope.buttons = buttons;
    
    $scope.cancel = function() {
        $mdDialog.cancel();
    };

    $scope.confirm = function() {
        $mdDialog.hide();
    };

    // Focus on the cancel button to avoid accidental confirmation.  Do this
    // is as a zero-delay timeout so that the focus happens after Angular does
    // its refresh, else we'll try to focus the button before it exists.
    $timeout(function() { document.getElementById('cancel').focus(); });
});
