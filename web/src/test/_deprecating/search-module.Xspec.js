import ApiService from '@/service/api'
import UserService from '@/service/user'
import FileService from '@/service/file'
import SearchService from '@/service/search'

// Test files services.
describe('FilesService', function () {
  beforeAll((done) => {
    ApiService.init()

    Object.assign(this, {
      login: 'testErr',
      email: 'testErr@example.org',
      password: 'testErr1',

      groupId: 'aTestGroup',
      groupName: '',

      file: new Blob(['<a id="a"><p>I have seen...all good people...and I like it.</p><b id="b">And then, and that!</b></a>'],
        { type: 'text/html' }),
      fileId: '[TEST_SETUP]',
      fileName: 'blob',
      find: 'and',

      directoryId: 'KKwUc_hx',
      directoryName: 'aTestDirectory'
    })
    done()
  })

  // Create a user, ignore duplicate user name error.
  beforeAll((done) => {
    UserService.create({
      email: this.email,
      name: this.login,
      pass: this.password
    }).then(() => {
      done()
    }).catch((error) => {
      error.message.indexOf('409') > 0 ||
        fail(`Test Suite Setup: ${error.message}`)
      done()
    })
  })

  // Login
  beforeAll((done) => {
    UserService.login({
      name: this.login,
      pass: this.password
    }).then((data) => {
      done()
    }).catch((error) => {
      fail(`Test Suite Setup: userlogin ${error.message}`)
      done()
    })
  })

  // Create some files.
  //  - OR -
  // Store a file ID from an existing file.
  beforeAll((done) => {
    let find = this.find
    let file = new Blob(['<a id="a"><p>I have seen...all good people...and I like it.</p><b id="b">And then, and that!</b></a>'],
      { type: 'text/html' })
    SearchService.user({ find })
      .then(({ data }) => {
        this.fileId = data.matched[0].file.id
        done()
      })
      .catch(({ message }) => {
        FileService.create({ file }).then(({ data }) => {
          this.fileId = data
          FileService.create({ file }).then(({ data }) => {
            done()
          })
        }).catch(({ message }) => {
          fail(message)
          done()
        })
      })
  })

  /*
   * Begin Tests ***
   */
  describe('search', () => {
    let find = ''
    beforeAll(() => {
      find = this.find
    })

    it('e2e: returns [matched] property', (done) => {
      SearchService.search({ find })
        .then(({ data }) => {
          expect(Object.keys(data)).toEqual([ 'matched' ])
          done()
        })
        .catch(error => {
          fail(error.message)
          done()
        })
    })
  })

  xdescribe('listFiles', () => {
    it('e2e: returns [folders], and [files] properties', (done) => {
    // return ApiService.query(`record${context ? '/' + context : ''}`, { gid })

      FileService.list({})
        .then(({ data }) => {
          expect(Object.keys(data)).toEqual(['folders', 'files'])
          done()
        }).catch(({ message }) => {
          fail(message)
          done()
        })
    })
  })
})
