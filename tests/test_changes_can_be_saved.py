## Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.

import json
import time
from integration import DoccerTestCase


class SaveDocumentTest(DoccerTestCase):

    def test_documents_can_be_saved(self):
        name = 'Updating doccer doc'
        content = 'Here is some content for the doc\nNoice\n'
        sha1hash = 'e2c3a5c4570451054eef098a11b9ad10257c8cbe'
        self.api('/save', 'POST', data=dict(name=name))
        data = dict(name=name, content=content)

        self.api('/save', 'POST', data=data)

        with open('%s/content/%s' % (self.PATH, sha1hash,)) as f:
            self.assertEqual(f.read(), content)
        with open('accounts/default/%s' % (name,)) as f:
            contents = f.read()
            print contents
            first_save, second_save, _ = contents.split('\n')
            self.assertEqual(
                first_save.split()[0],
                '13cb80d32f4ce1c32da7fab8df5d6ae629a665a8'
            )
            self.assertAlmostEqual(
                int(first_save.split()[1]),
                int(time.time()), delta=10
            )
            self.assertEqual(
                second_save.split()[0],
                '39fa695bc6ad917122ca44da375466515fbfc0ed'
            )
            self.assertAlmostEqual(
                int(second_save.split()[1]),
                int(time.time()), delta=10
            )

    def test_latest_version_of_doc_content_is_returned(self):
        name = 'Latest doccer doc'
        content = 'Here is some content for the doc\nNoice\n'
        self.api('/save', 'POST', data=dict(name=name))
        data = dict(name=name, content=content)
        self.api('/save', 'POST', data=data)

        response = json.loads(self.api('/doc/Latest+doccer+doc', 'GET'))

        self.assertEqual(response['content'], content)

    def test_blobs_can_be_retrieved(self):
        name = 'Latest doccer doc'
        content = 'Here is some content for the doc\nNoice\n'
        data = dict(name=name, content=content)
        history = json.loads(self.api('/save', 'POST', data=data))['history']
        response = self.api('/blob/%s' % (history[0]['hash'],), 'GET')
        content_hash = response.split()[-1]

        response = self.api('/blob/%s' % (content_hash,), 'GET')

        self.assertIn(content, response)

    def test_document_history_is_available(self):
        name = 'History doccer doc'
        content = 'Here is some content for the doc\nNoice\n'
        self.api('/save', 'POST', data=dict(name=name))
        data = dict(name=name, content=content)
        self.api('/save', 'POST', data=data)
        data2 = dict(name=name, content=content + ' Some more content')
        self.api('/save', 'POST', data=data2)

        data = json.loads(self.api('/doc/History+doccer+doc', 'GET'))
        print data

        self.assertEqual(data2['content'], data['content'])
        self.assertEqual(data['name'], 'History doccer doc')
        self.assertEqual(data['encoded'], 'History+doccer+doc')
        self.assertEqual(data['history'][0]['hash'],
            '13cb80d32f4ce1c32da7fab8df5d6ae629a665a8')
        self.assertEqual(
            data['history'][1]['hash'],
            '39fa695bc6ad917122ca44da375466515fbfc0ed')
        self.assertEqual(
            data['history'][2]['hash'],
            '40514713c89982b21804e120dda2f2673ca5c333')
        self.assertAlmostEqual(data['history'][0]['ts'],
                int(time.time()), delta=10)
        self.assertAlmostEqual(data['history'][1]['ts'],
                int(time.time()), delta=10)
        self.assertAlmostEqual(data['history'][2]['ts'],
                int(time.time()), delta=10)
