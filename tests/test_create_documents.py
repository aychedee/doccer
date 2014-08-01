## Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.

from integration import DoccerTestCase


class CreateDocumentTest(DoccerTestCase):

    def test_new_documents_can_be_created(self):
        data = dict(name='My Latest Notes')
        self.api('/save', 'post', data=data)

        with open('%s/content/%s' % (self.PATH, self.EMPTY_HASH)) as f:
            self.assertEqual(f.read(), '')
