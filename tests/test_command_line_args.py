#!/usr/bin/env python3

from subprocess import Popen
import requests

from integration import DoccerTestCase


class CommandLineArgs(DoccerTestCase):

    @classmethod
    def setUpClass(cls):
        Popen('%s/doccer %s %s &' % (cls.PATH, '0.0.0.0', 9988), shell=True)

    def test_port_and_address_can_be_passed_as_arg(self):
        response = requests.get('http://127.0.0.1:9988/')

        self.assertIn('Doccer', response.content)
