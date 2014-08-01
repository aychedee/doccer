## Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.

from integration import DoccerTestCase


class HelloWorldTest(DoccerTestCase):

   def test_server_returns_valid_http_response(self):
        response = self.api('/', 'GET')

        self.assertIn('Doccer', response)
