function Doccer() {

    var Promise = require('promise');     
    var http = require('http');
    'use strict';

        this.api = function(endpoint, method, callback, data) {

            var req = http.request(
                {
                    method: method,
                    path: endpoint 
                }, function(response) {});

            req.onreadystatechange = function() {
                if (req.readyState == 4 ) {
                   if(req.status == 200){
                       var resp = JSON.parse(xhr.responseText);
                       console.log(resp);
                       callback(resp);
                   }
                   else if(req.status == 400) {
                      console.log('There was an error 400')
                   }
                   else {
                       console.log('xhr.status was ' + req.status)
                   }
                }
            }
            if (method.toLowerCase() === 'post') {
                req.setHeader("Content-type", "application/x-www-form-urlencoded");
                req.end(data);
            } else {
                req.end();
            }
        }

        this.newDoc = function(name, callback) {
            this.api(
                '/new', 'POST', callback,
                'name=' + encodeURIComponent(name)
            )
            return new Promise(function(resolve, deny) {});
        }

        this.save = function(name, content, callback) {
            this.api(
                '/save', 'POST', callback,
                'name=' + encodeURIComponent(name) + '&' + 'content=' + encodeURIComponent(content)
            )
        }

        this.listDocs = function(callback) {
            this.api('/docs', 'GET', callback)
        }

        this.getDoc = function(name, cHash, callback) {
            console.log('Callback', callback, 'Name', name)
            var cb = function(callback, resp) {
                this.doc = new this.Doc(resp);
                return callback(resp);
            }
            this.api('/doc/' + encodeURIComponent(name) + '/' + cHash, 'GET', cb.bind(this, callback))
        }

        this.getBlob = function(hash, callback) {
            this.api('/blob/' + hash, callback)
        }

        function Doc(doc, doccer) {
            this.name = doc.name;
            this.doccer = doccer;
            this.content = doc.content;
            this.encoded = doc.encoded;
            this.history = [];
            this.position = doc.history.length;
            for (var i=0; i < doc.history.length; i++) {
                this.history.push({
                    date: new Date(doc.history[i].ts),
                    hash: doc.history[i].hash,
                    content: ''
                });
            }
        }

        this.Doc = Doc;

}

module.exports = Doccer;
