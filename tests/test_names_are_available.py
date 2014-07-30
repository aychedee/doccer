#!/usr/bin/env python3

import json

from integration import DoccerTestCase


class ListDocsTest(DoccerTestCase):

    def test_documents_can_be_listed(self):
        name1 = 'A doccer document'
        name2 = 'My second+'
        self.api('/new', 'POST', data=dict(name=name1), status_code=301)
        self.api('/new', 'POST', data=dict(name=name2), status_code=301)

        content = self.api('/docs/', 'GET')
        print content

        data = json.loads(content)
        print data
        self.assertEqual(data[0]['name'], 'A doccer document')
        self.assertEqual(data[0]['encoded'], 'A+doccer+document')
        self.assertEqual(
            data[0]['hash'], 'd62e44d32f5f21f83c67ce12d7c2ef0b29ca80d5')
        self.assertEqual(data[1]['name'], 'My second+')
        self.assertEqual(data[1]['encoded'], 'My+second%2B')
        self.assertEqual(
            data[1]['hash'], '20f5a7ef0502abca78de471c47beca182907ccff')
