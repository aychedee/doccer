#!/usr/bin/env python3

import requests

from integration import DoccerTestCase


class CreateDocumentTest(DoccerTestCase):

    def test_new_documents_can_be_created(self):
        sha1hash = '8d6ac39986ccd929d7cc1825efb0faa841a46e0a'
        data = dict(name='My Latest Notes')
        self.api(
            '/new', 'post', data=data, status_code=301,
            headers=dict(location='/doc/My+Latest+Notes')
        )

        with open('%s/content/%s' % (self.PATH, sha1hash,)) as f:
            self.assertEqual(f.read(), 'My Latest Notes\n===============\n')

    def test_server_returns_valid_http_response(self):
        response = requests.get('http://127.0.0.1:9999/')

        self.assertIn('Doccer', response.content)
