
var viewDB = angular.module("viewDB", ["ngRoute",'ngAnimate', 'ngSanitize', 'ui.bootstrap', 'switcher'])

viewDB.config(function($interpolateProvider, $routeProvider, $locationProvider) {
    $interpolateProvider.startSymbol('{[{').endSymbol('}]}');
    $locationProvider.html5Mode(true);
    $routeProvider
        .when('/', {
            templateUrl: '/static/html/search.html',
            controller: 'searchCtrl',
            resolve: {
                viewDB: function ($q) {
                    var defer = $q.defer();
                    defer.resolve();
                    return defer.promise;
                }
            }
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
    $scope.currentPage = 0
    //elastic search
    $scope.query = {
        "from" : 0, "size" : 10,
            "query": {
                "bool": {
                    "should": {
                        "multi_match": {
                            "query": $scope.phrase,
                            "fields": ["name_and_requisits", "government_choice", "reg_date", "id"]
                        }
                    }
                }
            }
    }
    $scope.addPhrase = function() {
        $scope.query.query.bool.should.multi_match.query = $scope.phrase
    }
    $scope.addTerm = function (field, data) {
        if (!$scope.query.query.bool.filter) {
            $scope.query.query.bool.filter = {
                "bool": {
                    "should": []
                }
            }
        }
        var obj = {}
        var arr = $scope.query.query.bool.filter.bool.should
        obj[field] = data
        if (arr.length == 0) {
            arr.push({"term": obj})
        } else {
            for (var i = 0 ; i < arr.length; i++){
                if (arr[i].term.hasOwnProperty(field)) {
                    arr[i].term[field] = data
                    break
                } else if(i === arr.length - 1 && !arr[i].term.hasOwnProperty(field)) {
                    arr.push({"term": obj})
                }
            }
        }

    }
    // $scope.logpage = function (page) {
    //     console.log(page)
    // }
    $scope.doSearch = function() {
        $http.post("http://192.168.99.100:9200/views/_search",$scope.query).then(function (response) {
            $scope.results = response.data;
        })
    }

    // get gov names and ids
    $http.get("/api/govs").then(function(response) {
        $scope.govs = response.data;
    })

});

viewDB.controller("viewCtrl", function($scope, $http, $sce, $location) {

    var id = $location.absUrl().slice(33)

    $scope.years = [2003,2004,2005,2006,2007,2008,2009,2011,2012,2013,2014,2015,2016,2017,2018,2019,2020,2021,2022,2023,2024]

    $http.get("/api/v/"+ id).then(function(response) {
        $scope.track = response.data;
    })

    $http.get("/api/govs").then(function(response) {
        $scope.govs = response.data;
    })

// start saving value
    $scope.saveData = function(msg, name, newValue, oldValue) {
        console.log(msg, name, id, "|", newValue, oldValue);
        $http.post("/api/v/"+ id, {
            name: name,
            data: newValue,
            id:id
        }).then(function(response) {
            $scope.track.pl[name] = newValue;
        })
    };
    $scope.log = function(msg, name, newValue, oldValue) {
        console.log(msg, name, "|",newValue, oldValue);
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