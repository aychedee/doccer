#!/usr/bin/env python3

import json

from integration import DoccerTestCase


class ListDocsTest(DoccerTestCase):

    def test_documents_can_be_listed(self):
        name1 = 'A doccer document'
        name2 = 'My second'
        self.api('/doc/new', 'POST', data=dict(name=name1), status_code=301)
        self.api('/doc/new', 'POST', data=dict(name=name2), status_code=301)

        content = self.api('/docs/', 'GET')

        data = json.loads(content)
        self.assertEqual(data[0]['name'], 'A doccer document')
        self.assertEqual(
            data[0]['hash'], 'd62e44d32f5f21f83c67ce12d7c2ef0b29ca80d5')
        self.assertEqual(data[1]['name'], 'My second')
        self.assertEqual(
            data[1]['hash'], 'c0dff7df435d556b01c1d8bae31540a7c109ce61')
