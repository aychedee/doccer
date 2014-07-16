#!/usr/bin/env python3

import requests

from integration import DoccerTestCase


class CreateDocumentTest(DoccerTestCase):

    def test_new_documents_can_be_created(self):
        sha1hash = '8d6ac39986ccd929d7cc1825efb0faa841a46e0a'
        data = dict(name='My Latest Notes')
        response = requests.post(
            'http://127.0.0.1:9999/doc/new', data=data, allow_redirects=False)

        self.assertEqual(response.status_code, 301)
        self.assertEqual(
            response.headers['location'],
            '/doc/%s' % (sha1hash,)
        )
        with open('%s/content/%s.md' % (self.PATH, sha1hash,)) as f:
            self.assertEqual(f.read(), 'My Latest Notes\n===============\n')

    def test_server_returns_valid_http_response(self):
        response = requests.get('http://127.0.0.1:9999/')

        self.assertIn('Doccer', response.content)
