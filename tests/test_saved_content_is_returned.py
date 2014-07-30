#!/usr/bin/env python3

import json

from integration import DoccerTestCase


class CreateDocumentTest(DoccerTestCase):

    def test_new_documents_can_be_created(self):
        name = 'Just the default content'
        self.api('/new', 'POST', data=dict(name=name), status_code=301)
        content = self.api('/docs/', 'GET')
        data = json.loads(content)

        content = self.api('/blob/%s' % (data[0]['hash'],), 'GET')

        self.assertIn(name, content)
