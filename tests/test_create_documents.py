## Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.

import requests

from integration import DoccerTestCase


class CreateDocumentTest(DoccerTestCase):

    def test_new_documents_can_be_created(self):
        data = dict(name='My Latest Notes')
        self.api(
            '/new', 'post', data=data, status_code=301,
            headers=dict(location='/doc/My+Latest+Notes')
        )

        with open('%s/content/%s' % (self.PATH, self.EMPTY_HASH)) as f:
            self.assertEqual(f.read(), '')

    def test_server_returns_valid_http_response(self):
        response = requests.get('http://127.0.0.1:9999/')

        self.assertIn('Doccer', response.content)
