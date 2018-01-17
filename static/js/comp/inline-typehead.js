 viewDB.directive('inlineTypehead', function ($timeout) {
        return {
            scope: {
                model: '=inlineTypehead',
                governs: '=typeGovs',
                handleSave: '&onSave',
                handleCancel: '&onCancel'
            },
            link: function (scope, elm, attr) {
                var previousValue;
                $timeout(function () {
                        scope.data = scope.governs[scope.model]["name"]
                    },
                    100);
                scope.edit = function () {
                    scope.editMode = true;
                    previousValue = scope.model;

                    $timeout(function () {
                        elm.find('input')[0].focus();
                    }, 0, false);
                };
                scope.save = function () {
                    scope.editMode = false;
                    scope.handleSave({value: String(scope.model)});
                };
                scope.cancel = function () {
                    scope.editMode = false;
                    scope.model = previousValue;
                    scope.handleCancel({value: scope.model});
                };
                scope.formatLabel = function(gmodel) {
                    for (var i=0; i< scope.governs.length; i++) {
                        if (gmodel.id === scope.governs[i].id) {
                            scope.model = scope.governs[i].id
                        }
                    }
                };
            },
            templateUrl: '/static/html/inline/inline-typehead.html'
        };
    });
