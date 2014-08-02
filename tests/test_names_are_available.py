## Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.
# -*- coding: utf8

import json

from integration import DoccerTestCase


class ListDocsTest(DoccerTestCase):

    def test_documents_can_be_listed(self):
        name1 = 'A doccer document'
        name2 = 'My second+'
        self.api('/save', 'POST', data=dict(name=name1))
        self.api('/save', 'POST', data=dict(name=name2))

        data = json.loads(self.api('/docs/', 'GET'))

        self.assertEqual(data[0]['name'], 'A doccer document')
        self.assertEqual(data[0]['encoded'], 'A+doccer+document')
        self.assertEqual(
            data[0]['history'][0]['hash'],
           '13cb80d32f4ce1c32da7fab8df5d6ae629a665a8')
        self.assertEqual(data[1]['name'], 'My second+')
        self.assertEqual(data[1]['encoded'], 'My+second%2B')
        self.assertEqual(
            data[1]['history'][0]['hash'],
           '13cb80d32f4ce1c32da7fab8df5d6ae629a665a8')
