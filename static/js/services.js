
viewDB.service('authService', function ($http) {
        /*jshint validthis: true */
        const baseURL = 'http://localhost:8888/';

            this.login = function(user) {
                return $http({
                    method: 'POST',
                    url: baseURL + 'u/login',
                    data: user,
                    headers: {'Content-Type': 'application/json'}
                });
            };
            this.logout = function(user) {
                return $http({
                    method: 'POST',
                    url: baseURL + 'u/logout',
                    data: user,
                    headers: {'Content-Type': 'application/json'}
                });
            };
            this.register =  function(user) {
                return $http({
                    method: 'POST',
                    url: baseURL + 'u/register',
                    data: user,
                    headers: {'Content-Type': 'application/json'}
                });
            };
            this.ensureAuthenticated = function(token) {
                return $http({
                    method: 'GET',
                    url: baseURL + 'api/cabinet',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: 'Bearer ' + token
                    }
                });
            };
});

viewDB.service('trackingService', function ($http) {
    /*jshint validthis: true */
    const baseURL = 'http://localhost:8888/';
    this.getGovs = function () {
        return $http({
            method: 'GET',
            url: baseURL + 'api/govs',
        });
    }
});

viewDB.service('fileUploadService', function ($http, $q) {

    this.uploadFileToUrl = function (file, uploadUrl, docId, token) {
        //FormData, object of key/value pair for form fields and values
        var fileFormData = new FormData();
        fileFormData.append('file', file);
        fileFormData.append('doc_id', docId);
        var deffered = $q.defer();
        $http.post(uploadUrl, fileFormData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined},
            Authorization: 'Bearer ' + token

        }).then(function successCallback(response) {
            deffered.resolve(response);

        }, function errorCallback(response) {
            deffered.reject(response);
        });

        return deffered.promise;
    }
});
