## Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.
# -*- coding: utf8

import json

from integration import DoccerTestCase


class UnicodeTest(DoccerTestCase):

    def test_document_list_with_unicode(self):
        name1 = u'☠   ⚓'
        name2 = u'⚓⚓⚓'
        self.api('/new', 'POST', data=dict(name=name1), status_code=301)
        self.api('/new', 'POST', data=dict(name=name2), status_code=301)

        data = json.loads(self.api('/docs/', 'GET'))

        self.assertEqual(data[0]['name'], name1)
        self.assertEqual(data[0]['encoded'], u'%E2%98%A0+++%E2%9A%93')
        self.assertEqual(
            data[0]['history'][0]['hash'], self.EMPTY_HASH)
        self.assertEqual(data[1]['name'], name2)
        self.assertEqual(data[1]['encoded'], u'%E2%9A%93%E2%9A%93%E2%9A%93')
        self.assertEqual(
            data[1]['history'][0]['hash'], self.EMPTY_HASH)
