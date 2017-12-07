viewDB.directive('inlineTypehead', function($timeout) {
    return {
        scope: {
            model: '=inlineTypehead',
            govs: '=typeGovs',
            handleSave: '&onSave',
            handleCancel: '&onCancel'
        },
        link: function(scope, elm, attr) {
            var previousValue;
            $timeout(function () {
                scope.data = scope.govs[scope.model]["name"]
            },
                    100);
            scope.edit = function() {
                scope.editMode = true;
                previousValue = scope.model;

                $timeout(function() {
                    elm.find('input')[0].focus();
                }, 0, false);
            };
            scope.save = function() {
                scope.editMode = false;
                scope.handleSave({value: String(scope.model.Id)});
            };
            scope.cancel = function() {
                scope.editMode = false;
                scope.model = previousValue;
                scope.handleCancel({value: scope.model});
            };
        },
        templateUrl: '/static/html/inline-typehead.html'
    };
});