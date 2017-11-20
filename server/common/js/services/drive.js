host.factory('drive', function ($q, $http, $filter, bug, schema) {
    'use strict';

    function __chainReject(p, error) {
        p.reject(error);
    }

    function chainReject(p) {
        return __chainReject.bind(null, p);
    }

    var __client = $q.defer();
    gapi.load('client:auth2', function() {
        gapi.client.init({
            clientId: '383959393597-qgmja5j8us2gr06m4soca79qsgt6ku3c.apps.googleusercontent.com',
            fetch_basic_profile: false,
            discoveryDocs: ['https://www.googleapis.com/discovery/v1/apis/drive/v3/rest'],
            scope: 'https://www.googleapis.com/auth/drive.metadata.readonly https://www.googleapis.com/auth/drive.file',
        }).then(
            function() {
                if (gapi.auth2.getAuthInstance().isSignedIn.get()) {
                    __client.resolve();
                } else if (localStorage.getItem('drive') != $filter('date')(new Date(), 'yyyy-MM-dd')) {
                    localStorage.setItem('drive', $filter('date')(new Date(), 'yyyy-MM-dd'));
                    gapi.auth2.getAuthInstance().signIn().then(
                        function() {
                            __client.resolve();
                        },
                        chainReject(__client)
                    );
                }
            },
            chainReject(__client)
        );
    });

    function getFileId(name) {
        var i = $q.defer();
        __client.promise.then(
            function() {
                gapi.client.drive.files.list({
                    q: "mimeType='application/json' and trashed=false and name='{0}'".format(name),
                    fields: 'files(id, name)',
                    spaces: 'drive',
                }).then(
                    function(response) {
                        var files = response.result.files;
                        var id = files && files.length > 0 ? files[0].id : undefined;
                        if (id) {
                            i.resolve(id);
                        } else {
                            gapi.client.drive.files.create({
                                mimeType: 'application/json',
                                name: name,
                                fields: 'id'
                            }).then(
                                function(response) {
                                    i.resolve(response.result.id);
                                },
                                chainReject(i)
                            );
                        }
                    },
                    chainReject(i)
                );
            },
            chainReject(i)
        );
        return i.promise;
    }

    function patchFile(id, ops) {
        return gapi.client.request({
            path: '/upload/drive/v3/files/' + id,
            method: 'PATCH',
            params: { uploadType: 'media' },
            body: JSON.stringify(ops)
        });
    }

    function getFile(id) {
        var i = $q.defer();
        gapi.client.drive.files.get({
            fileId: id,
            mimeType: 'application/json',
            alt: 'media',
        }).then(
            function(response) {
                i.resolve(response.body ? JSON.parse(response.body) : {});
            },
            chainReject(i)
        );
        return i.promise;
    }

    function fakePatch(id, lambda) {
        // Google Drive is fucking retarded.  It supports, and
        // basically requires, the PATCH method and goes on and
        // about PATCH semantics, but doesn't actually support
        // patching anything.  So, manually fudge PATCH via GET
        // and PUT (well, PATCH, as Google dictates).
        var i = $q.defer();
        getFile(id).then(
            function(file) {
                patchFile(id, lambda(file)).then(
                    function() {
                        i.resolve(file);
                    },
                    chainReject(i)
                );
            },
            chainReject(i)
        );
        return i.promise;
    }

    return {
        get: function(name, initData) {
            var i = $q.defer();
            getFileId(name).then(
                function(id) {
                    getFile(id).then(
                        function(file) {
                            if (_.isEmpty(file)) {
                                file = initData();
                                patchFile(id, file).then(
                                    function() {
                                        i.resolve({ id: id, data: file });
                                    },
                                    chainReject(i)
                                );
                            } else {
                                i.resolve({ id: id, data: file });
                            }
                        },
                        chainReject(i)
                    );
                },
                chainReject(i)
            );
            return i.promise;
        },
        patch: function(id, pending) {
            return fakePatch(id, function(file) {
                _.each(pending, function(patches, index) {
                    file[index] = file[index] || {};
                    _.each(patches, function(patch, key) {
                        if (patch.add) {
                            file[index][key] = patch.val;
                        } else if (file[index].hasOwnProperty(key)) {
                            delete file[index][key];
                        }
                    });
                });
                return file;
            });
        },
    };
});