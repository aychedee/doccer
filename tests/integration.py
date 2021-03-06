## Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.

from unittest import TestCase
import requests
from subprocess import check_output, Popen
import os
from os.path import abspath, dirname
import time
import glob

PATH = abspath(dirname(dirname(__file__)))

class DoccerTestCase(TestCase):

    PATH = PATH
    DEFAULT_HOST = '127.0.0.1'
    DEFAULT_PORT = 4121
    doccer_location = '%s/doccer' % (PATH,)
    EMPTY_HASH = 'da39a3ee5e6b4b0d3255bfef95601890afd80709'

    @classmethod
    def setUpClass(cls):
        cls.SERVER_URL = cls.start_doccer()

    @classmethod
    def tearDownClass(cls):
        cls.stop_doccer()
        files = glob.glob(cls.PATH + '/accounts/default/*')
        for f in files:
            os.unlink(f)

    def api(self, endpoint, method, data={}, status_code=200, headers={}):
        response = getattr(requests, method.lower())(
            '%s%s' % (self.SERVER_URL, endpoint),
            data=data,
            allow_redirects=False
        )
        if headers:
            for key in headers.keys():
                self.assertEqual(response.headers[key], headers[key])
        self.assertEqual(response.status_code, status_code)

        return response.content

    @classmethod
    def doccer_running(cls):
        return cls.doccer_location in check_output(
            'ps auxf | grep doccer', shell=True)

    @classmethod
    def run_doccer(self, host=None, port=None):
        if host or port:
            command = '%s %s %s &' % (self.doccer_location, host, port)
        else:
            command = self.doccer_location + ' &'
        Popen(command, shell=True)

    @classmethod
    def start_doccer(cls, host=None, port=None):
        cls.run_doccer(host=host, port=port)
        while not cls.doccer_running():
            time.sleep(0.01)
        if host or port:
            return 'http://%s:%s' % (host, port)
        else:
            return 'http://%s:%s' % (cls.DEFAULT_HOST, cls.DEFAULT_PORT)

    @classmethod
    def stop_doccer(cls):
        Popen(['killall', 'doccer'])
        while cls.doccer_running():
            time.sleep(0.01)
