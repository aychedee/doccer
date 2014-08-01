Doccer
======

Doccer is a document storage server. It is intended to be used with a web based
front end. It can be installed on your local machine and used as a private note
store or it can be installed on a server and used by many many people.

The main point of difference between doccer and anything else? It keeps a
complete history of changes to any document you store in it. Every single save
point is available and can be traversed.


Usage
=====

Doccer has two concepts: **Documents** and **blobs**. A document has a name,
and a history or timeseries of blobs. When you save or update the content of a
document it appends a new blob to that history.

Blobs contain content. So to visually replay all the changes to your document
over time you would fetch each blob from the document history in turn and
display them to the user.

Documents are owned by accounts, but if the account system is not turned on
then all documents are stored in the default account.

The JSON representation of a document looks like this:

    {
        "name": "My document",
        "encoded": My%20Document",
        "history": [
            {
                "hash": "da39a3ee5e6b4b0d3255bfef95601890afd80709",
                "ts": "2015-11-10T23:00:00Z"
            },
            {
                "hash": "e12e169cda777c5a751000e30a46ac77e7d2ffdd",
                "ts": "2015-11-10T23:15:00Z"
            },
        ],
        "content": "Doc content"
    }

To step back to the first time the document was saved you would retrieve the
oldest blob by making a GET request to

    /blob/da39a3ee5e6b4b0d3255bfef95601890afd80709.


JavaScript API
--------------

The JavaScript API makes implementing a web based frontend very simple.

    doccer = new Doccer()

    // creating a new document
    doccer.save(
        'Shopping list',
        'Some content',
        function(resp) {console.log(resp);}
    );

    // updating the content of the same document
    doccer.save(
        'Shopping list',
        'Some content, plus some more content',
        function(resp) {console.log(resp);}
    );

    // getting a list of the documents for your account
    doccer.listDocs(function(resp) {console.log(resp);});
    // The content attribute of all these docs will be empty

    // retreiving an individual document
    doccer.getDoc('Shopping list', function(resp) {console.log(resp);});

    // fetching the content blob for a document's 4th modification
    doccer.getBlob(doc.history[3].hash, function(resp) {console.log(resp);});
