var interceptor = function ($q, $location) {
    return {
        request: function (config) {//req
            if (!localStorage.getItem('token')) {
                localStorage.setItem('token', 'unknown');
            }
            return config;
        },

        response: function (result) {//res
            return result;
        },

        responseError: function (rejection) {
            if (rejection.status == 401 && rejection.config.method == 'POST') {
                $location.url('/u/login');
            }
            return $q.reject(rejection);
        }
    }
};
viewDB.config(function($interpolateProvider, $routeProvider, $locationProvider, $httpProvider) {


    $httpProvider.interceptors.push(interceptor);
    $interpolateProvider.startSymbol('{[{').endSymbol('}]}');
    $locationProvider.html5Mode(true);
    $routeProvider
        .when('/', {
            title: 'Пошук відстежень',
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
            title: 'Якісні показники ЦОВВ',
            templateUrl: '/static/html/ratings.html',
            controller: 'ratingsCtrl',
            resolve: {
                viewDB: function ($q) {
                    var defer = $q.defer();
                    defer.resolve();
                    return defer.promise;
                }
            }
        })
        .when('/u/register', {
            title: 'Реєстрація',
            templateUrl: '/static/html/auth/auth.register.view.html',
            controller: 'authRegisterController',
        })
        .when('/u/login', {
            title: 'Вхід',
            templateUrl: '/static/html/auth/auth.login.view.html',
            controller: 'authLoginController',
        })
        .when('/u/cabinet', {
            title: 'Кабінет користувача',
            templateUrl: '/static/html/auth/auth.cabinet.view.html',
            controller: 'userCabinetController',
        })
        .when('/track/id/:trackId', {
            templateUrl: '/static/html/view.html',
            controller: 'viewCtrl',
        })
        .when('/track/id/:trackId/edit', {
            templateUrl: '/static/html/edit.html',
            controller: 'editCtrl',
        })
        .when('/track/create', {
            title: 'Додати відстеження',
            templateUrl: '/static/html/create.html',
            controller: 'createCtrl',
        })
        .when('/govs/edit', {
            title: 'Редагувати назви державних органів',
            templateUrl: '/static/html/govs.html',
            controller: 'govsEditController',
        });

});
viewDB.run(['$rootScope', '$route', function($rootScope, $route) {
    $rootScope.$on('$routeChangeSuccess', function() {
        document.title = $route.current.title;
    });
}]);