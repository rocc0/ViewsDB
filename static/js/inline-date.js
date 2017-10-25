viewDB.directive('inlineDate', function($timeout) {
    return {
        scope: {
            model: '=inlineDate',
            handleSave: '&onSave',
            handleCancel: '&onCancel'
        },
        link: function(scope, elm, attr) {
            var previousValue;

            scope.today = function() {
                scope.model = new Date();
            };
            scope.today();

            scope.clear = function() {
                scope.model = null;
            };
            scope.dateOptions = {
                formatYear: 'yy',
                maxDate: new Date(2028, 5, 22),
                minDate: new Date(1991, 5, 22),
                startingDay: 1
            };
            scope.popup = {
                opened: false
            };

            scope.open = function() {
                scope.popup.opened = true;
            };

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
        templateUrl: '/static/html/inline-date.html'
    };
});