var viewDB = angular.module("viewDB", []);
viewDB.config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('{[{').endSymbol('}]}');
});

viewDB.controller("viewDBCtrl", function($scope, $http) {

    // editor
    $http.get('/api/v/2').then(function(response) {
        $scope.viewData = response.data.pl;
    })

    $scope.getTemplate = function (contact) {
        if (contact.k === $scope.viewData.selected) return 'edit';
        else return 'display';
    };

    $scope.editContact = function (contact) {
        $scope.viewData.selected = angular.copy(contact);
    };

    $scope.saveContact = function (idx) {
        console.log("Saving contact");
        $scope.viewData.contacts[idx] = angular.copy($scope.viewData.selected);
        $scope.reset();
    };

    $scope.reset = function () {
        $scope.viewData.selected = {};
    };
    // end editor

    $http.get('/api/fields').then(function(response) {
            $scope.fieldNames = response.data.fields;
        })


    $scope.searchSyntax = function() {
        $http.post('/api/search', {
            "size": 10,
            "explain": true,
            "highlight":{},
            "query": {
                "boost": 1.0,
                "query": $scope.syntax,
            }
        }).
        then(function(response) {
            $scope.processResults(response.data);
        })
    };

    $scope.roundTook = function(took) {
        if (took < 1000 * 1000) {
            return "less than 1ms";
        } else if (took < 1000 * 1000 * 1000) {
            return "" + Math.round(took / (1000*1000)) + "ms";
        } else {
            roundMs = Math.round(took / (1000*1000));
            return "" + roundMs/1000 + "s";
        }
    };

    $scope.roundScore = function(score) {
        return Math.round(score*1000)/1000;
    };

    $scope.expl = function(explanation) {
        rv = "" + $scope.roundScore(explanation.value) + " - " + explanation.message;
        rv = rv + "<ul>";
        for(var i in explanation.children) {
            child = explanation.children[i];
            rv = rv + "<li>" + $scope.expl(child) + "</li>";
        }
        rv = rv + "</ul>";
        return rv;
    };

    $scope.processResults = function(data) {
        $scope.errorMessage = null;
        $scope.results = data;
        for(var i in $scope.results.hits) {
            hit = $scope.results.hits[i];
            hit.roundedScore = $scope.roundScore(hit.score);
            hit.explanationString = $scope.expl(hit.explanation);
            hit.explanationStringSafe = $sce.trustAsHtml(hit.explanationString);
            for(var ff in hit.fragments) {
                fragments = hit.fragments[ff];
                newFragments = [];
                for(var ffi in fragments) {
                    fragment = fragments[ffi];
                    safeFragment = $sce.trustAsHtml(fragment);
                    newFragments.push(safeFragment);
                }
                hit.fragments[ff] = newFragments;
            }
        }
        $scope.results.roundTook = $scope.roundTook(data.took);
    };
});


