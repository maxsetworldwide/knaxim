describe('Knaxim', function () {
  beforeAll(function (done) {
    this.url = '/api'
    this.knc = new Knaxim(this.url)

    Object.assign(this, {
      username: 'testErr',
      email: 'testErr@example.org',
      password: 'testErr1',

      groupId: 'aTestGroup',
      groupName: '',

      file: new Blob(['<a id="a"><b id="b">hey!</b></a>'],
        { type: 'text/html' }),

      fileId: 'ArIIevng', // StorageId
      fileName: 'blob',

      directoryId: 'KKwUc_hx',
      directoryName: 'aTestDirectory'
    })
    /*
    this.knc.createUser(this.username, this.email, this.password).then(res => {
      this.knc.login(this.username, this.password).then(res => {
       if (res.modified == 'User Created' || res.modified == 'Name Taken') {
         done()
       }
     })
    }).catch(e => {
      fail('Test Suite Setup: createUser')
      done()
    })
    */
  })

  beforeEach(function (done) {
    this.knc.login(this.username, this.password).then(res => {
      done()
    }).catch(e => {
      fail('Test Suite Setup: createUser')
      done()
    })
  })

  // Begin Tests ***
  // ***************
  // ***************
  describe('downloadpath', function () {
    it('returns the correct path', function () {
      expect(this.knc.downloadpath('5'))
        .toEqual(`${this.url}/file/5/download`)
    })
  })

  // * is tested during test suite setup.  If additonal vectors or tests are
  // needed, enable this test group and add them here.
  xdescribe('createUser', function () {
    let username = 'localtestErr1'
    let email = 'localtestErr1@example.org'
    let password = 'localtestErr1'

    it('gets something', function (done) {
      this.knc.createUser(username, email, password).then(res => {
        expect(res.modified).toBe('User Created')
        // TODO: Delete User
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('createAdmin', function () {
    let username = 'testErrAdmin'
    let email = 'testErrAdmin@example.org'
    let password = 'testErrAdmin1'
    let adminkey = 'b9!5v789DEp$5Yqw@9h' // copied from Knaxim server config

    it('creates a user admin', function (done) {
      // console.log('TODO: The server appears to be adding a newline char')

      this.knc.createAdmin(username, email, password, adminkey).then(res => {
        expect(res.modified).toBe('User Created')
        done()
      }).catch(e => {
        expect(e.message).toBe('Conflict')
        done()
      })
    })
  })

  describe('currentUser', function () {
    beforeEach(function () {
    })

    it('gets the current user', function (done) {
      // console.log('TODO: Does the server return content-type application/json')

      this.knc.currentUser().then(res => {
        expect(res.modified.id).toBe(this.username)
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  // * is tested during test suite setup.  If additonal vectors or tests are
  // needed, enable this test group and add them here.
  xdescribe('login', function () {
    it('gets something', function (done) {
      this.knc.login(this.username, this.password).then(res => {
        expect(res.modified).toBe('logged in')
        done()
      })
    })
  })

  // This functionality was disabled in the API.
  xdescribe('_confirmEmail', function () {
    beforeEach(function () {
    })

    it('gets something', function (done) {
      this.knc._confirmEmail().then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('signout', function () {
    beforeEach(function () {
    })

    it('logs user out', function (done) {
      this.knc.signout().then(res => {
        expect(res.modified).toBe('Signed Out')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('getUserData', function () {
    beforeEach(function () {
    })

    it('gets user storage info', function (done) {
      // console.log('TODO: Does the server return content-type application/json')

      this.knc.getUserData().then(res => {
        expect(res.modified.totaldata).toBe(52428800)
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('searchUser', function () {
    let find = 'macroMeString'

    it('gets content summary', function (done) {
      console.log('TODO: Does the server return content-type application/json')

      this.knc.searchUser(find).then(res => {
        expect(res.modified.size).toBe(0)
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('changePassword', function () {
    let newPass = 'testErr2'
    let password

    beforeAll(function () {
      password = this.password
    })

    it('changes the password', function (done) {
      this.knc.changePassword(password, newPass).then(res => {
        this.knc.changePassword(newPass, password).then(res => {
          done()
        })
      }).catch(e => {
        fail('Test cleanup may have failed when trying to reset the password')
        done()
      })
    })
  })

  describe('userProfile', function () {
    let username

    beforeAll(function () {
      username = this.username
    })

    it('gets a user profile', function (done) {
      console.log('TODO: Does the server return content-type application/json')

      this.knc.userProfile().then(res => {
        expect(res.modified.user.name).toBe(username)
        done()
      }).catch(e => {
        done()
      })
    })
  })

  describe('getPermissions', function () {
    let type = 'group' // {'dir' | 'group' | 'file'} type
    let id

    beforeEach(function () {
      id = this.groupId
    })

    it('returns Not Found', function (done) {
      console.log('TODO: Does the server return content-type application/json')

      this.knc.getPermissions(type, id).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Not Found')
        done()
      })
    })
  })

  describe('setPermission', function () {
    let type = 'group' // {'dir' | 'group' | 'file'}
    let record
    let user = 'macroMeString'
    let key = 'macroMeString'
    let value = true

    beforeEach(function () {
      record = this.groupId
    })

    it('returns Not Found', function (done) {
      console.log('TODO: Mock resource on Test server?')

      this.knc.setPermission(type, record, user, key, value).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Not Found')
        done()
      })
    })
  })

  describe('setPermissionPublic', function () {
    let type = 'macroMeString'
    let record = 'group' // {'dir' | 'group' | 'file'}
    let key
    let value = true

    beforeEach(function () {
      key = this.groupId
    })

    it('sets public permissions', function (done) {
      this.knc.setPermissionPublic(type, record, key, value).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Not Found')
        done()
      })
    })
  })

  describe('getRecords', function () {
    let key = 'testErr'
    let group

    beforeEach(function () {
      group = this.groupId
    })

    it('returns Not Found', function (done) {
      this.knc.getRecords(key, group).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Not Found')
        done()
      })
    })
  })

  describe('setRecordName', function () {
    let id = 'macroMeString'
    let name = 'macroMeString'

    it('gets something', function (done) {
      // TODO: Set Record Name of file
      // TODO: Test Cleanup reSet record name
      console.log('TODO: Is invalid data handled correctly w/500?')

      this.knc.setRecordName(id, name).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        fail(e.message)
        done()
      })
    })
  })

  describe('createGroup', function () {
    let name = 'aTestGroupTmp'
    let maker

    beforeAll(function () {
      maker = this.username
    })

    it('creates a group', function (done) {
      console.log('TODO: Is invalid data handled correctly w/500?')
      console.log(' - creating duplicate group returns 500?')

      this.knc.createGroup(name, maker).then(res => {
        expect(res.modified).toBe('Group Created')
        fail('TODO: deleteGroup when the method becomes available.')
        done()
      }).catch(e => {
        fail(e.message)
        done()
      })
    })
  })

  describe('getGroups', function () {
    let group

    beforeAll(function () {
      group = this.groupId
    })

    it('gets group metadata', function (done) {
      console.log('TODO: Does the server return content-type application/json')

      this.knc.getGroups(group).then(res => {
        expect(res.modified.member).toBe('Array')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })

    it('gets a groups metadata for the current user', function (done) {
      this.knc.getGroups().then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('searchGroup', function () {
    let find = 'test'
    let id

    beforeAll(function () {
      id = this.groupId
    })

    it('finds content that matches a search term', function (done) {
      console.log('TODO: Does the server return content-type application/json')

      this.knc.searchGroup(find, id).then(res => {
        expect(res.modified.size).toBe(0)
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('searchPublic', function () {
    let find = 'aSearchString'

    it('gets a list of documents', function (done) {
      this.knc.searchPublic(find).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('getGroupMembers', function () {
    let group

    beforeAll(function () {
      group = this.groupId
    })

    it('gets member info', function (done) {
      console.log('TODO: Does the server return content-type application/json')
      console.log('TODO: Mock resource on Test server?')

      this.knc.getGroupMembers(group).then(res => {
        expect(res.modified.gid).toBe(group)
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('changeGroupMember', function () {
    let group
    let user = 'testUserA'
    let add = true

    beforeAll(function () {
      group = this.groupId
    })

    it('adds current user to their own group??', function (done) {
      this.knc.changeGroupMember(group, this.username, add).then(res => {
        expect(res.modified).toBe('updated members')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })

    it('removes current user from their own group??', function (done) {
      this.knc.changeGroupMember(group, this.username, !add).then(res => {
        expect(res.modified).toBe('updated members')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })

    xit('adds a user', function (done) {
    })
    xit('removes a user', function (done) {
    })
  })

  describe('createDirectory', function () {
    let name = 'aTestDirectoryTmp'
    let owner
    let content = ['abc', '123'] // Array of file id's

    beforeAll(function () {
      owner = this.username
    })

    it('creates a directory', function (done) {
      console.log('Creating a directory with fake file ids ignores file ids??')
      console.log('Creating a w/ duplicate name creates a new directory??')

      this.knc.createDirectory(name, owner, content).then(res => {
        expect(res.modified.name).toBe('aTestDirectory')
        fail('TODO: Delete directory')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('createDynDir', function () {
    let name = 'aTestDynDir'
    let owner
    let contexttype = 'dir' // 'dir' | 'group' | 'user'
    let contextid = 'aTestDirectory'
    let search = 'aTestTerm'

    beforeAll(function () {
      owner = this.username
    })

    it('creates a dynamic dir', function (done) {
      console.log('TODO: Is invalid data handled correctly w/500?')
      console.log('TODO: Mock resource on Test server?')

      this.knc.createDynDir(name, owner, contexttype, contextid, search).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('refreshDynDir', function () {
    let id = 'aTestDynDir'

    it('Updates document contents', function (done) {
      console.log('TODO: Is invalid data handled correctly w/500?')
      console.log('TODO: Mock resource on Test server?')

      this.knc.refreshDynDir(id).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('getDirectory', function () {
    let id

    beforeAll(function () {
      id = this.directoryId
    })

    it('gets a directory', function (done) {
      console.log('TODO: Does the server return content-type application/json')
      console.log('TODO: non-existing IDs cause 500 error?')

      this.knc.getDirectory(id).then(res => {
        expect(res.modified.id).toBe(id)
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('editDirectory', function () {
    let id
    let fileid
    let add = true

    beforeAll(function () {
      id = this.directoryId,
      fileid = this.fileId
    })

    it('adds a file', function (done) {
      this.knc.editDirectory(id, fileid, add).then(res => {
        console.log('TODO: Is invalid data handled correctly w/500?')

        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })

    xit('removes a file', function (done) {
    })
  })

  describe('searchDirectory', function () {
    let id
    let find = 'testTerm'

    beforeAll(function () {
      id = this.directoryId
    })

    it('gets a list of file ids', function (done) {
      this.knc.searchDirectory(id, find).then(res => {
        console.log('TODO: Is invalid data handled correctly w/500?')

        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('deleteDirectory', function () {
    let name = 'aTestDirectoryToDelete'
    let owner
    let content = ['abc', '123'] // Array of file id's
    let id

    beforeEach(function (done) {
      owner = this.username

      this.knc.createDirectory(name, owner, content).then(res => {
        id = res.modified.id
        done()
      }).catch(e => {
        if (e.message == 'Content Exists?') {
          done()
        }
        fail('Test Setup: cr$eateDirectory')
        done()
      })
    })

    it('deletes a directory', function (done) {
      this.knc.deleteDirectory('randO').then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('createFile', function () {
    let owner
    let folder //      folder = 'o_5KZHw4',
    let aFileParts = ['<a id="a"><b id="b">hey!</b></a>']
    let file = new Blob(aFileParts, { type: 'text/html' })

    beforeAll(function (done) {
      owner = this.username
    })

    it('creates a file', function (done) {
      console.log('TODO: Does the server return content-type application/json')
      console.log('TODO: Name ends up being "blob", is another param needed??')

      this.knc.createFile(owner, file, folder).then(res => {
        expect(res.modified.name).toBe('blob')
        done()
        /*
        // TODO: Cleanup might need to run in afterAll...
        this.knc.deleteFile(owner, res.modified.id).then(res => {
          done()
        }).catch(e => {
          fail('Test Setup: deleteFile()')
          done()
        })
        */
      }).catch(e => {
        done()
      })
    })
  })

  describe('trackWebPage', function () {
    let url = 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf'
    let owner
    let folder

    beforeEach(function () {
      owner = this.username
    })

    it('creates a web resource', function (done) {
      this.knc.trackWebPage(url, owner, folder).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('refreshWebPage', function () {
    let url = 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf'
    let owner
    let folder
    let id

    beforeEach(function (done) {
      owner = this.username

      this.knc.trackWebPage(url, owner, folder).then(res => {
        id = res.modified.id
        done()
      }).catch(e => {
        done()
      })
    })

    it('refreshes a web resource', function (done) {
      this.knc.refreshWebPage(id).then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('getFile', function () {
    let id

    beforeAll(function () {
      id = this.fileId
    })

    it('gets a file record', function (done) {
      this.knc.getFile(id).then(res => {
        expect(res.modified.id).toBe(id)
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('getFileSlice', function () {
    let id
    let start = 0
    let end = 1

    beforeAll(function () {
      id = this.fieldId
    })

    it('gets a range of sentences', function (done) {
      this.knc.getFileSlice(id, start, end).then(res => {
        expect(res.modified.size).toBe(1)
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('searchFileSlice', function () {
    let id
    let start = 0
    let end = 1
    let find = 'hey'

    beforeAll(function () {
      id = this.fieldId
    })

    it('returns slices by search term', function (done) {
      this.knc.searchFileSlice(id, start, end, find).then(res => {
        expect(res.modified.size).toBe(1)
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })

  describe('deleteFile', function () {
    let id
    let owner
    let folder
    let aFileParts = ['<a id="a"><b id="b">hey!</b></a>']
    let file = new Blob(aFileParts, { type: 'text/html' })

    beforeAll(function (done) {
      owner = this.username

      this.knc.createFile(owner, file, folder).then(res => {
        id = res.modified.id
        done()
      }).catch(e => {
        fail('Test Setup: createFile()')
        done()
      })
    })

    it('deletes a file', function (done) {
      this.knc.deleteFile(id).then(res => {
        expect(res.modified).toBe('Record Removed')
        done()
      }).catch(e => {
        fail(e.message)
        done()
      })
    })
  })

  // Disabled in api for some reason.
  xdescribe('_deleteGroup', function () {
    beforeEach(function () {
    })

    it('gets something', function (done) {
      this.knc._deleteGroup().then(res => {
        expect(res.modified).toBe('something')
        done()
      }).catch(e => {
        expect(e.message).toBe('Unauthorized')
        done()
      })
    })
  })
})
