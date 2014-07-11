#!/usr/bin/env python3

from unittest import TestCase
from subprocess import Popen
from os.path import abspath, dirname
import time

PATH = abspath(dirname(dirname(__file__)))

class DoccerTestCase(TestCase):

    PATH = PATH

    @classmethod
    def setUpClass(cls):
        Popen('%s/doccer &' % (PATH,), shell=True)
        time.sleep(3)
        Popen('ps auxf | grep doccer', shell=True)

    @classmethod
    def tearDownClass(cls):
        Popen(['killall', 'doccer'])
