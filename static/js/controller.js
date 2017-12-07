
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
        .when('/track/id/:trackId', {
            templateUrl: '/static/html/view.html',
            controller: 'viewCtrl',
        })
        .when('/track/id/:trackId/edit', {
            templateUrl: '/static/html/edit.html',
            controller: 'editCtrl',
        })
        .when('/track/create', {
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
                            "fields": ["requisits", "gov_choice", "reg_date"]
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

    $scope.doSearch = function() {
        $http.post("http://192.168.99.100:9200/tracking/_search",$scope.query).then(function (response) {
            $scope.results = response.data;
        })
    }

    // get gov names and ids
    $http.get("/api/govs").then(function(response) {
        $scope.govs = response.data;
    })

});


viewDB.controller("viewCtrl", function ($scope, $http, $location,$routeParams) {
    $scope.params = $routeParams
    var docId = $scope.params.trackId
    $http.get("/api/v/"+ docId).then(function(response) {
        $scope.track = response.data;
        for (k in $scope.track.pl) {
            if ($scope.track.pl[k] == '') {
                $scope.track.pl[k] = 'Інформація відсутня'
            }
        }
        $http.get("/api/govs").then(function(response) {
            $scope.govs = response.data;
        });
    });



    $scope.printDiv = function(divName) {
        var printContents = document.getElementById(divName).innerHTML;
        var originalContents = document.body.innerHTML;
        document.body.innerHTML = printContents;
        window.print();
        document.body.innerHTML = originalContents;
    }
})


viewDB.controller("editCtrl", function($scope, $http, $sce, $location,$routeParams) {
    $scope.params = $routeParams
    var docId = $scope.params.trackId

    var period = {
        "trace_id": docId,
        "term_per": "",
        "res_per_bool": 0,
        "res_per_year": 0,
        "sign_per": "",
        "publ_per": "",
        "give_per": "",
        "res_per": "не заповнено",
        "concl_per": "не заповнено",
        "cp_bool": 0
    };
    $scope.addPeriod = function(data,reg_date) {
        console.log(!data,reg_date)
        if (!data || data.length == 0) {
            data = []
            reg_date = new Date(reg_date)
            new_term = angular.copy(reg_date)
            new_term.setFullYear(reg_date.getFullYear()+4)
            var month = parseInt(reg_date.getMonth())
            if (month < 12) {
                month += 1
            }
            period.term_per = new_term.getFullYear() + "-" + month + "-" + new_term.getDate()
            $scope.savePeriod(period)
            data.push(period)
            $scope.track.pr = data
            delete period.pid
        } else {
            date = data[data.length-1].term_per
            date = new Date(date)
            period = angular.copy(period)
            new_term = angular.copy(date)
            new_term.setFullYear(date.getFullYear()+3)
            var month = parseInt(date.getMonth())
            if (month < 12) {
                month += 1
            }
            period.term_per = new_term.getFullYear() + "-" + month + "-" + new_term.getDate()
            $scope.savePeriod(period)
            data.push(period)
            $scope.track.pr = data
            delete period.pid
        }

    }
    $scope.savePeriod = function(data) {
        $http.post("/api/create", data
        ).then(function(response) {
            console.log(response.data.id)
            period.pid = response.data.id
        })
    };
    $scope.removePeriod = function (index, model) {
        $http.post("/api/delete",{
            "item_id": model[index].pid,
            "tbl_name": "p"
        }).then(function (response) {
            model.splice(index,1)
            console.log(response.data)
        })

    }
    $scope.years = [2003,2004,2005,2006,2007,2008,2009,2011,
        2012,2013,2014,2015,2016,2017,2018,2019,2020,2021,2022,2023,2024];

    $http.get("/api/v/"+ docId).then(function(response) {
        $scope.track = response.data;
    });

    $http.get("/api/govs").then(function(response) {
        $scope.govs = response.data;
    });

// start saving value
    $scope.saveData = function(name,newValue,oldValue) {
        console.log(name, docId, "|", newValue, oldValue);
        $http.post("/api/v/"+ docId, {
            type: 0,
            name: name,
            data: newValue,
            id:parseInt(docId)
        }).then(function(response) {
            $scope.track.pl[name] = new newValue;
        })
    };

    $scope.saveDataPer = function(name, pid, value) {
        $http.post("/api/v/"+ pid, {
            type: 1,
            name: name,
            data: value,
            id:pid
        }).then(function(response) {
            $scope.track.pr[name] = value;
        })
    };
    $scope.log = function(name, newValue) {
        console.log(name, "|",newValue);
    };
    // end saving value

    //image
    $scope.name = "Select Files to Upload";
    $scope.images = [];
    $scope.display = $scope.images[$scope.images.length - 1];
    $scope.setImage = function (ix) {
        $scope.display = $scope.images[ix];
    }
    $scope.clearAll = function () {
        $scope.display = '';
        $scope.images = [];
    }
    $scope.upload = function (obj) {
        var elem = obj.target || obj.srcElement;
        for (i = 0; i < elem.files.length; i++) {
            var file = elem.files[i];
            var reader = new FileReader();

            reader.onload = function (e) {
                $scope.images.push(e.target.result);
                $scope.display = e.target.result;
                $scope.$apply();
            }
            reader.readAsDataURL(file);
        }
    }
    //image

});

viewDB.controller("createCtrl", function($scope, $http) {
    $http.get("/api/govs").then(function(response) {
        $scope.govs = response.data;
    })

    $scope.addTrack = function () {
        $http.post("/api/create",$scope.fdata).then(function(response) {
            $scope.pdata.trace_id = response.data.id;
            if ($scope.postPeriod == 0) {
                $http.post("/api/create",$scope.pdata).then(function(response) {
                    window.location.replace("/track/id/"+$scope.pdata.trace_id);
                })
            } else {
                window.location.replace("/track/id/"+$scope.pdata.trace_id);
            }

        })
    }
    $scope.dateConvert = function (model) {
        var value = model
        if (new Date(model) != 'Invalid Date') {
            var month = parseInt(model.getMonth())
            if (month < 12) {
                month += 1
            }
            value = model.getFullYear() + "-" + month + "-" + model.getDate()
        }
        return value
    }


});

viewDB.controller("ratingsCtrl", function($scope, $http) {
    $scope.printDiv = function(divName) {
        var printContents = document.getElementById(divName).innerHTML;
        var originalContents = document.body.innerHTML;

        document.body.innerHTML = printContents;

        window.print();

        document.body.innerHTML = originalContents;
    }
    $http.get("/api/ratings" ).then(function(response) {
        $scope.ratings = response.data
    })
});