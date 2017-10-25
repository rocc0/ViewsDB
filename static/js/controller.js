var viewDB = angular.module("viewDB", ["ngRoute",'ngAnimate', 'ngSanitize', 'ui.bootstrap', 'ui.toggle'])

viewDB.config(function($interpolateProvider, $routeProvider, $locationProvider) {
    $interpolateProvider.startSymbol('{[{').endSymbol('}]}');
    $locationProvider.html5Mode(true);
    $routeProvider
        .when('/', {
            templateUrl: '/static/html/search.html',
            controller: 'viewDBCtrl',
        })
        .when('/ratings', {
            templateUrl: '/static/html/ratings.html',
            controller: 'ratingsCtrl',
            resolve: {
                viewDB: function ($q) {
                    var defer = $q.defer();
                    defer.resolve();
                    return defer.promise;
                }
            }})
        .when('/views/view/:id', {
            templateUrl: '/static/html/view.html',
            controller: 'viewCtrl',
        })
        .when('/views/create', {
            templateUrl: '/static/html/create.html',
            controller: 'createCtrl',
        });
});

viewDB.controller("viewDBCtrl", function($scope, $http, $sce) {
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

viewDB.controller("viewCtrl", function($scope, $http, $sce, $location) {

    var id = $location.absUrl().slice(33)

    $scope.years = [2003,2004,2005,2006,2007,2008,2009,2011,2012,2013,2014,2015,2016,2017,2018,2019,2020,2021,2022,2023,2024]

    $http.get("/api/v/"+ id).then(function(response) {
        $scope.viewData = response.data;
    })

    $http.get("/api/govs").then(function(response) {
        $scope.govs = response.data;
    })

// start saving value
    $scope.saveData = function(msg, name, data) {
        console.log(msg, name, data, id);
        $http.post("/api/v/"+ id, {
            name: name,
            data: data,
            id:id
        }).then(function(response) {
            $scope.viewData.pl[name] = data;
        })
    };
    $scope.log = function(msg, value) {
        console.log(msg, value);
    };
    // end saving value
});


viewDB.controller("ratingsCtrl", function($scope, $http) {
    $scope.form_data = {};
    $http.get("/api/ratings").then(function(response) {
        $scope.ratings = response.data
    })
});

viewDB.controller("createCtrl", function($scope, $http) {
    "use strict";
    $http.get("/api/ratings").then(function(response) {
        $scope.ratings = response.data
    })
});