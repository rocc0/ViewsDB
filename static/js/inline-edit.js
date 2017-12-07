viewDB.directive('inlineEdit', function($timeout) {
    return {
        scope: {
            model: '=inlineEdit',
            handleSave: '&onSave',
            handleCancel: '&onCancel'
        },
        link: function(scope, elm, attr) {
            var previousValue;
            scope.date = function () {
                if (scope.model == "")  {
                    scope.model = "інф. відсутня"
                }
            }
            scope.date()
            scope.edit = function() {
                scope.editMode = true;
                previousValue = scope.model;
                $timeout(function() {
                    elm.find('input')[0].focus();
                }, 0, false);
            };
            scope.save = function() {
                scope.editMode = false;
                scope.handleSave({value: scope.model});
            };
            scope.cancel = function() {
                scope.editMode = false;
                scope.model = previousValue;
                scope.handleCancel({value: scope.model});
            };
        },
        templateUrl: '/static/html/inline-edit.html'
    };
});