// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
