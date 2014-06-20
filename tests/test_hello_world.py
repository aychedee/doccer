#!/usr/bin/env python3

from unittest import TestCase
import os
import requests



class HelloWorldTest(TestCase):

    def test_server_returns_valid_http_response(self):

        os.system('/home/hansel/Coding/gowiki/notes 127.0.0.1 9999 &')

        response = requests.get('http://127.0.0.1:9999/')

        self.assertEqual(response.content, 'Hello, World')


