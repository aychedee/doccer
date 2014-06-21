#!/usr/bin/env python3

from unittest import TestCase
from subprocess import Popen
import requests
from os.path import abspath, dirname

PATH = abspath(dirname(dirname(__file__)))
print PATH

class DoccerTestCase(TestCase):

    @classmethod
    def setUpClass(cls):
        Popen('%s/doccer &' % (PATH,), shell=True)

    @classmethod
    def tearDownClass(cls):
        Popen(['killall', 'doccer'])


class HelloWorldTest(DoccerTestCase):



    def test_server_returns_valid_http_response(self):
        response = requests.get('http://127.0.0.1:9999/')

        self.assertIn('Doccer', response.content)


class CreateDocumentTest(DoccerTestCase):

    def test_new_documents_can_be_created(self):
        pass


class CommandLineArgs(DoccerTestCase):

    @classmethod
    def setUpClass(cls):
        Popen('%s/doccer %s %s &' % (PATH, '0.0.0.0', 9988), shell=True)

    def test_port_and_address_can_be_passed_as_arg(self):
        response = requests.get('http://127.0.0.1:9988/')

        self.assertIn('Doccer', response.content)
