#!/usr/bin/env python3

import time
from integration import DoccerTestCase


class SaveDocumentTest(DoccerTestCase):

    def test_documents_can_be_saved(self):
        name = 'Updating doccer doc'
        content = 'Here is some content for the doc\nNoice\n'
        sha1hash = 'e2c3a5c4570451054eef098a11b9ad10257c8cbe'
        self.api('/new', 'POST', data=dict(name=name), status_code=301)
        data = dict(name=name, content=content)

        response = self.api('/save', 'POST', data=data)

        self.assertEqual(response, sha1hash)
        with open('%s/content/%s' % (self.PATH, sha1hash,)) as f:
            self.assertEqual(f.read(), content)
        with open('accounts/default/%s' % (name,)) as f:
            contents = f.read()
            print contents
            first_save, second_save, _ = contents.split('\n')
            self.assertEqual(
                first_save.split()[0],
                self.EMPTY_HASH
            )
            self.assertAlmostEqual(
                int(first_save.split()[1]),
                int(time.time()), delta=10
            )
            self.assertEqual(
                second_save.split()[0],
                'e2c3a5c4570451054eef098a11b9ad10257c8cbe'
            )
            self.assertAlmostEqual(
                int(second_save.split()[1]),
                int(time.time()), delta=10
            )

    def test_latest_version_of_doc_content_is_returned(self):
        name = 'Latest doccer doc'
        content = 'Here is some content for the doc\nNoice\n'
        self.api('/new', 'POST', data=dict(name=name), status_code=301)
        data = dict(name=name, content=content)
        self.api('/save', 'POST', data=data)

        response = self.api('/doc/Latest+doccer+doc', 'GET')

        self.assertIn(content, response)

    def test_blobs_can_be_retrieved(self):
        name = 'Latest doccer doc'
        content = 'Here is some content for the doc\nNoice\n'
        self.api('/new', 'POST', data=dict(name=name), status_code=301)
        data = dict(name=name, content=content)
        response = self.api('/save', 'POST', data=data)

        response = self.api('/blob/%s' % (response,), 'GET')

        self.assertIn(content, response)
