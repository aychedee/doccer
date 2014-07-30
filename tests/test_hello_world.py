## Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.

import requests

from integration import DoccerTestCase


class HelloWorldTest(DoccerTestCase):

   def test_server_returns_valid_http_response(self):
        response = requests.get('http://127.0.0.1:9999/')

        self.assertIn('Doccer', response.content)
