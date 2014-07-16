#!/usr/bin/env python3

import requests

from integration import DoccerTestCase


class CommandLineArgs(DoccerTestCase):

    @classmethod
    def setUpClass(cls):
        cls.start_doccer(host='0.0.0.0', port=9988)

    def test_port_and_address_can_be_passed_as_arg(self):
        response = requests.get('http://127.0.0.1:9988/')

        self.assertIn('Doccer', response.content)
