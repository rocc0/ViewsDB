var viewDB = angular.module("viewDB", ["ngRoute",'ngAnimate', 'ngSanitize', 'ui.bootstrap', 'ui.toggle'])

viewDB.config(function($interpolateProvider, $routeProvider, $locationProvider) {
    $interpolateProvider.startSymbol('{[{').endSymbol('}]}');
    $locationProvider.html5Mode(true);
    $routeProvider
        .when('/', {
            templateUrl: '/static/html/search.html',
            controller: 'searchCtrl',
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

viewDB.controller("searchCtrl", function($scope, $http) {
    $scope.years = [2003,2004,2005,2006,2007,2008,2009,2011,2012,2013,2014,2015,2016,2017,2018,2019,2020,2021,2022,2023,2024]
    this.filter = function () {
        if ($scope.search_id != null) {
            return "filter" : {
                "term": {
                    "id": $scope.search_id
                }
            }
        } else {

        }
    }
    $scope.searchSyntax = function () {
        $http.post("http://192.168.99.100:9200/views/_search",{
            "query":{
                "bool": {
                    "should" : {
                        "multi_match" : {
                            "query": $scope.phrase,
                            "fields": [ "name_and_requisits", "government_choice", "reg_date", "id" ]
                        }
                    },
                }}
            }).then(function (response) {
            $scope.results = response.data;
        })
    }
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