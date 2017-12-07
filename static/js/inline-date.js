viewDB.directive('inlineDate', function($timeout) {
    return {
        scope: {
            model: '=inlineDate',
            handleSave: '&onSave',
            handleCancel: '&onCancel'
        },
        link: function(scope, elm, attr) {
            var previousValue;

            scope.clear = function() {
                scope.model = null;
            };
            if (scope.model != '') {
                scope.model = new Date(scope.model)
            }

            scope.edit = function() {
                scope.editMode = true;
                if (scope.model != '') {
                    scope.model = new Date(scope.model)
                }
                previousValue = scope.model;

                $timeout(function() {
                    elm.find('input')[0].focus();
                }, 0, false);
            };
            scope.save = function() {
                var month = parseInt(scope.model.getMonth())
                if (month < 12) {
                    month += 1
                }
                date_value = scope.model.getFullYear() + "-" + month + "-" + scope.model.getDate()
                scope.editMode = false;
                scope.handleSave({value: date_value});
            };
            scope.cancel = function() {
                scope.editMode = false;
                scope.model = previousValue;
                scope.handleCancel({value: scope.model});
            };
        },
        templateUrl: '/static/html/inline-date.html'
    };
});