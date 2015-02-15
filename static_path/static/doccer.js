function Doccer() {
    'use strict';

    this.api = function(endpoint, method, callback, data) {
        var xhr;

        if (window.XMLHttpRequest) {
            // code for IE7+, Firefox, Chrome, Opera, Safari
            xhr = new XMLHttpRequest();
        } else {
            // code for IE6, IE5
            xhr = new ActiveXObject("Microsoft.XMLHTTP");
        }

        xhr.onreadystatechange = function() {
            if (xhr.readyState == 4 ) {
               if(xhr.status == 200){
                   var resp = JSON.parse(xhr.responseText);
                   console.log(resp);
                   callback(resp);
               }
               else if(xhr.status == 400) {
                  console.log('There was an error 400')
               }
               else {
                   console.log('xhr.status was ' + xhr.status)
               }
            }
        }
        xhr.open(method, endpoint, true);
        if (method.toLowerCase() === 'post') {
            xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
            xhr.send(data);
        } else {
            xhr.send();
        }
    }

    this.newDoc = function(name, callback) {
        this.api(
            '/new', 'POST', callback,
            'name=' + encodeURIComponent(name)
        )
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
