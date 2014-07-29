#!/usr/bin/env python3

import json

from integration import DoccerTestCase


class CreateDocumentTest(DoccerTestCase):

    def test_new_documents_can_be_created(self):
        name1 = 'A doccer document'
        name2 = 'My second'
        self.api('/doc/new', 'POST', data=dict(name=name1))
        self.api('/doc/new', 'POST', data=dict(name=name2))

        content = self.api('/docs/', 'GET')

        self.assertEqual(status_code, 200)
        data = json.loads(content)
        self.assertEqual(data[0]['name'], 'A doccer document')
        self.assertEqual(
            data[0]['hash'], 'd62e44d32f5f21f83c67ce12d7c2ef0b29ca80d5')
        self.assertEqual(data[1]['name'], 'My Latest Notes')
        self.assertEqual(
            data[1]['hash'], '8d6ac39986ccd929d7cc1825efb0faa841a46e0a')
        self.assertEqual(data[2]['name'], 'My second')
        self.assertEqual(
            data[2]['hash'], 'c0dff7df435d556b01c1d8bae31540a7c109ce61')
