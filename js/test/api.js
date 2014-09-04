'user strict';

var chai = require("chai");
var chaiAsPromised = require("chai-as-promised");
chai.use(chaiAsPromised);
var assert = chai.assert;
        


var Doccer = require("../doccer.js");
var Promise = require('promise');     

describe('newDoc', function() {

    before(function() {
    });

    after(function() {
    });

    beforeEach(function() {
        this.doccer = new Doccer();
        this.doccer.api = function() {};
    });

    afterEach(function() {
    });

    it('should create new doc object from server response', function() {
        var docPromise = this.doccer.newDoc('New doc');

        assert.becomes(docPromise, 'New doc');

    });

    it('should return a promises object', function() {
        this.doccer.api = function() {};

        var docPromise = this.doccer.newDoc('Test document', 'initial content');

        assert.instanceOf(docPromise, Promise);

    });

});

describe('Doc', function() {

    beforeEach(function() {
        this.doccer = new Doccer();
        this.doccer.api = function() {};
    });

    it('should make new doc from json data and doccer instance', function() {
        var doc = new this.doccer.Doc({
            name: 'A name',
            content: 'Some content',
            encoded: 'encoding',
            history: [
                {
                    ts: '',
                    hash: ''
                },
                {
                    ts: '',
                    hash: ''
                },
                {
                    ts: '',
                    hash: ''
                },
                {
                    ts: '',
                    hash: ''
                }
            ]
        }, this.doccer);

        assert.instanceOf(doc, this.doccer.Doc)
        assert.equal(doc.doccer, this.doccer)
        assert.equal(doc.content, 'Some content')
        assert.equal(doc.name, 'A name')


    });

});
