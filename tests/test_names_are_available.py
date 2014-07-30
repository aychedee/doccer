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
            data[0]['hash'], self.EMPTY_HASH)
        self.assertEqual(data[1]['name'], 'My second+')
        self.assertEqual(data[1]['encoded'], 'My+second%2B')
        self.assertEqual(
            data[1]['hash'], self.EMPTY_HASH)
