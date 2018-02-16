 var viewDB = angular.module("viewDB", ["ngRoute", 'ngAnimate', 'ngSanitize',
        'ui.bootstrap', 'switcher', 'bootstrapLightbox'])

 viewDB.controller("searchCtrl", function ($scope, $http, trackingService) {
        $scope.currentPage = 0
        //elastic search
        $scope.query = {
            "from": 0, "size": 10,
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
        $scope.addPhrase = function () {
            $scope.query.query.bool.should.multi_match.query = $scope.phrase
        }

        $scope.addTerm = function (field, data) {
            if (!$scope.query.query.bool.filter) {
                $scope.query.query.bool.filter = {
                    "bool": {
                        "must": []
                    }
                }
            }
            var obj = {}
            var arr = $scope.query.query.bool.filter.bool.must
            obj[field] = data
            if (arr.length == 0) {
                arr.push({"term": obj})
            } else {
                for (var i = 0; i < arr.length; i++) {
                    if (arr[i].term.hasOwnProperty(field)) {
                        arr[i].term[field] = data
                        break
                    } else if (i === arr.length - 1 && !arr[i].term.hasOwnProperty(field)) {
                        arr.push({"term": obj})
                    }
                }
            }

        }

        $scope.doSearch = function () {
            $http.post("http://192.168.99.100:9200/tracking/_search", $scope.query).then(function (response) {
                $scope.results = response.data;
            })
        }

        // get gov names and ids
        trackingService.getGovs()
           .then(function(response) {
            $scope.govs = response.data.govs;
        })

    });

 viewDB.controller("viewCtrl", function ($scope, $http, $location, $routeParams,trackingService) {
        $scope.params = $routeParams
        $scope.editPath = $location.path() + "/edit"
        var docId = $scope.params.trackId
        $http.get("/api/v/" + docId).then(function (response) {
            $scope.track = response.data;
            for (var k in $scope.track.pl) {
                if ($scope.track.pl[k] == '') {
                    $scope.track.pl[k] = 'Інформація відсутня'
                }
            }
            trackingService.getGovs()
                .then(function(response) {
                    $scope.govs = response.data;
                })
        });


        $scope.printDiv = function (divName) {
            var printContents = document.getElementById(divName).innerHTML;
            var originalContents = document.body.innerHTML;
            document.body.innerHTML = printContents;
            window.print();
            document.body.innerHTML = originalContents;
        }
    })

 viewDB.controller("editCtrl", function ($scope, $http, $sce, $location, trackingService,
                                            $routeParams, fileUploadService, Lightbox,authService) {
         const token = localStorage.getItem('token');
         if (token) {
             authService.ensureAuthenticated(token)
                 .then(function(user) {
                 })
                 .catch(function(err) {
                     console.log(err)
                     $location.path('/u/login');
                 });

         }
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

        $http.get("/api/v/" + docId).then(function (response) {
            $scope.track = response.data;
        });



     //format label for typehead on select
     trackingService.getGovs()
         .then(function(response) {
             $scope.governs = response.data.govs;
         });

     //end format label for typehead on select


        $http.get("/api/img/" + docId).then(function (response) {
            $scope.images = response.data.images;
        })

        $scope.years = [2003, 2004, 2005, 2006, 2007, 2008, 2009, 2011,
            2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024];

        $scope.addPeriod = function (data, reg_date) {
            console.log(!data, reg_date)
            if (!data || data.length == 0) {
                data = []
                reg_date = new Date(reg_date)
                new_term = angular.copy(reg_date)
                new_term.setFullYear(reg_date.getFullYear() + 4)
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
                var date = data[data.length - 1].term_per
                date = new Date(date)
                period = angular.copy(period)
                new_term = angular.copy(date)
                new_term.setFullYear(date.getFullYear() + 3)
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

        $scope.savePeriod = function (data) {
            $http({
                method: 'POST',
                url: "/api/create",
                data: data,
                headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token}
            }).then(function (response) {
                console.log(response.data.id)
                period.pid = response.data.id
            })
        };
        $scope.removePeriod = function (index, model) {
            $http({
                method: 'POST',
                url: "/api/delete",
                data: {item_id: model[index].pid, tbl_name: "p"},
                headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token}
        }).then(function (response) {
                model.splice(index, 1)
                console.log(response.data)
            })

        }

// start saving value
        $scope.saveData = function (name, newValue, oldValue) {
            console.log(newValue)
            $http({
                method: 'POST',
                url:"/api/v/" + docId,
                data: {type: 0, name: name, data: newValue, id: parseInt(docId)},
                headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
                }
            }).then(function (response) {
                $scope.track.pl[name] = newValue;
            })
        };

        $scope.saveDataPer = function (name, pid, value) {
            $http({
                method:'POST',
                url:"/api/v/" + pid,
                data: {type: 1, name: name, data: value, id: pid},
                headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token}
            }).then(function (response) {
                $scope.track.pr[name] = value;
            })
        };
        $scope.log = function (name, newValue) {
            console.log(name, "|", newValue);
        };
        // end saving value
        // image

        $scope.openLightboxModal = function (index) {
            Lightbox.openModal($scope.images, index);
        };
        // remove image
        $scope.removeImage = function (index) {
            var photo_id = $scope.images[index].photo_id
            $http({
                method: 'POST',
                url:"/api/img/" + docId + "/delete",
                data: {photo_id: photo_id},
                headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token}
            }).then(function (response) {
                $scope.images.splice(index, 1);

            })
        }
        // upload image
        $scope.uploadFile = function (images) {
            var uploadUrl = "/api/upload",
                promise = fileUploadService.uploadFileToUrl(images, uploadUrl, docId, token);

            promise.then(function (response) {
                if ($scope.images == null) {
                    $scope.images = []
                }
                $scope.images.push({
                    "photo_id": response.data.photo_id,
                    "original": response.data.original,
                    "thumb": response.data.thumb
                })
            }, function () {
                $scope.serverResponse = 'An error has occurred';
            })
        };
        //image

    });

 viewDB.controller("createCtrl", function ($scope, $http,trackingService) {
        trackingService.getGovs()
         .then(function(response) {
             $scope.governs = response.data.govs;
         });

        $scope.addTrack = function () {
            $http({
                method: 'POST',
                url: "/api/create",
                data: $scope.fdata,
                headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token}
            }).then(function (response) {
                window.location.replace("/track/id/" + response.data.id);
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

 viewDB.controller("ratingsCtrl", function ($scope, $http) {
        $scope.printDiv = function (divName) {
            var printContents = document.getElementById(divName).innerHTML;
            var originalContents = document.body.innerHTML;

            document.body.innerHTML = printContents;

            window.print();

            document.body.innerHTML = originalContents;
        }
        $http.get("/api/ratings").then(function (response) {
            $scope.ratings = response.data
        })
     $scope.today = function() {
         $scope.dt = new Date();
     };
     $scope.today();

     $scope.clear = function() {
         $scope.dt = null;
     };

     $scope.inlineOptions = {
         customClass: getDayClass,
         minDate: new Date(),
         showWeeks: true
     };

     $scope.dateOptions = {
         dateDisabled: disabled,
         formatYear: 'yy',
         maxDate: new Date(2020, 5, 22),
         minDate: new Date(),
         startingDay: 1
     };

     // Disable weekend selection
     function disabled(data) {
         var date = data.date,
             mode = data.mode;
         return mode === 'day' && (date.getDay() === 0 || date.getDay() === 6);
     }

     $scope.toggleMin = function() {
         $scope.inlineOptions.minDate = $scope.inlineOptions.minDate ? null : new Date();
         $scope.dateOptions.minDate = $scope.inlineOptions.minDate;
     };

     $scope.toggleMin();

     $scope.open1 = function() {
         $scope.popup1.opened = true;
     };

     $scope.open2 = function() {
         $scope.popup2.opened = true;
     };

     $scope.setDate = function(year, month, day) {
         $scope.dt = new Date(year, month, day);
     };

     $scope.formats = ['dd-MMMM-yyyy', 'yyyy/MM/dd', 'dd.MM.yyyy', 'shortDate'];
     $scope.format = $scope.formats[0];
     $scope.altInputFormats = ['M!/d!/yyyy'];

     $scope.popup1 = {
         opened: false
     };

     $scope.popup2 = {
         opened: false
     };

     var tomorrow = new Date();
     tomorrow.setDate(tomorrow.getDate() + 1);
     var afterTomorrow = new Date();
     afterTomorrow.setDate(tomorrow.getDate() + 1);
     $scope.events = [
         {
             date: tomorrow,
             status: 'full'
         },
         {
             date: afterTomorrow,
             status: 'partially'
         }
     ];

     function getDayClass(data) {
         var date = data.date,
             mode = data.mode;
         if (mode === 'day') {
             var dayToCheck = new Date(date).setHours(0,0,0,0);

             for (var i = 0; i < $scope.events.length; i++) {
                 var currentDay = new Date($scope.events[i].date).setHours(0,0,0,0);

                 if (dayToCheck === currentDay) {
                     return $scope.events[i].status;
                 }
             }
         }

         return '';
     }
    });

 viewDB.controller("authLoginController", function ($scope, $timeout, $location, authService) {
     $scope.user = {};
     $scope.onLogin = function() {
         authService.login($scope.user)
             .then(function(user) {
                 localStorage.setItem('token', user.data.token);
                 $location.path('/u/cabinet');
             })
             .catch(function(err) {
                 console.log(err);
                 $scope.message = "Невірний логін або пароль, спробуйте ще раз"
                 $timeout(function() {
                     $scope.message = ""
                 }, 2000);
             });
     };
 });

 viewDB.controller("userCabinetController", function ($scope, $http, $location, $rootScope, authService) {
     $scope.isLoggedIn = false;
     $scope.changepass = false

     const token = localStorage.getItem('token');
     if (token) {
         authService.ensureAuthenticated(token)
             .then(function(user) {
                 if (user.status === 200) {
                     $scope.userdata = user.data.data;
                     $scope.isLoggedIn = true;
                     $rootScope.isLoggedIn2 = true;
                 }
             })
             .catch(function(err) {
                 console.log(err)
                 $location.path('/u/login');
             });

     }
     $scope.postCheck = function () {
         checkService.postCheck(token)
             .then(function (response) {
                 $scope.testdata = response.data
             })
             .catch(function (err) {
                 console.log(err)
             })
     }
     $scope.getCheck = function () {
         checkService.getCheck(token)
             .then(function (response) {
                 $scope.testdata = response.data
             })
             .catch(function (err) {
                 console.log(err)
             })
     }

     $scope.changeUserField = function (field, id, value) {
         console.log(field, id, value)
         $http({
             method: 'POST',
             url:"/api/edituser",
             data: {field: field, data: value, id: parseInt(id)},
             headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
             }
         }).then(function (response) {
         }).catch(function(err){
             console.log(err)
         });
     };

 });

 viewDB.controller("authRegisterController", function ($scope,authService) {

     $scope.user = {
         password: "",
         confirmPassword: ""
     };
     $scope.onRegister = function() {
         authService.register(vm.user)
             .then(function(response) {
                 $location.path('/status');
             })
             .catch(function(err) {
                 console.log(err);
             });
     };
 });

 viewDB.controller("menuController", function ($scope, $rootScope,authService) {
     $scope.isLoggedIn = false;
     $rootScope.isLoggedIn2 = false;
     const token = localStorage.getItem('token');
     if (token) {
         authService.ensureAuthenticated(token)
             .then(function(user) {
                 if (user.status === 200) {
                     $scope.isLoggedIn = true;
                     $rootScope.isLoggedIn2 = true;
                 }
             })
             .catch(function(err) {

                 console.log(err)
             });

     }
     $scope.onLogout = function() {
         localStorage.removeItem('token');
         $scope.isLoggedIn = false;
         $rootScope.isLoggedIn2 = false;
         $location.path('/u/login');
     };
 });

 viewDB.controller("govsEditController", function ($scope,$http,$location,authService,trackingService) {
     const token = localStorage.getItem('token');
     if (token) {
         authService.ensureAuthenticated(token)
             .then(function(user) {
             })
             .catch(function(err) {
                 console.log(err)
                 $location.path('/u/login');
             });

     }
     trackingService.getGovs()
         .then(function(response) {
             $scope.govs = response.data.govs;
         })

     $scope.changeUserField = function (id, value) {
         console.log(id, value)
         $http({
             method: 'POST',
             url:"/api/govs/edit",
             data: {id: parseInt(id), name: value},
             headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
             }
         }).then(function(response) {
         }).catch(function(err){
             console.log(err)
         });
     };
 });
