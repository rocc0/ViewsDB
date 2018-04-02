viewDB.filter('boolToWord', function() {
    return function(object) {
        var word = "Ні"
        if (object == 1) {
            word = "Так"
        }
        return word;
    }
});