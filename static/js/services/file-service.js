viewDB.service('fileUploadService', function ($http, $q) {

        this.uploadFileToUrl = function (file, uploadUrl, docId) {
            //FormData, object of key/value pair for form fields and values
            var fileFormData = new FormData();
            fileFormData.append('file', file);
            fileFormData.append('docid', docId);
            var deffered = $q.defer();
            $http.post(uploadUrl, fileFormData, {
                transformRequest: angular.identity,
                headers: {'Content-Type': undefined}

            }).then(function successCallback(response) {
                deffered.resolve(response);

            }, function errorCallback(response) {
                deffered.reject(response);
            });

            return deffered.promise;
        }
    });
