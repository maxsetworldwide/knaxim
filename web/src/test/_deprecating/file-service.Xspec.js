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
import FileService from '@/service/file'
import setDefaults from '../test/setup'

describe('FileService', function () {
  // Test Setup
  beforeAll(async () => {
    await setDefaults(this)
  })

  /* Begin Tests */
  describe('slice', () => {
    let fid = ''
    let start = ''
    let end = ''

    beforeAll(() => {
      fid = this.fileId
      start = 0
      end = 1
    })

    it('e2e: returns size, and [lines] properties', (done) => {
      FileService.slice({ fid, start, end }).then(({ data }) => {
        expect(Object.keys(data)).toEqual(['size', 'lines'])
        done()
      }).catch(({ message }) => {
        fail(message)
        done()
      })
    })
  })

  describe('search', () => {
    let fid = ''
    let start = ''
    let end = ''
    let find = ''

    beforeAll(() => {
      fid = this.fileId
      start = 0
      end = 1
      find = 'and'
    })

    it('e2e: returns size, and [lines] properties', (done) => {
      FileService.search({ fid, start, end, find }).then(({ data }) => {
        expect(Object.keys(data)).toEqual(['size', 'lines'])
        done()
      }).catch(({ message }) => {
        fail(message)
        done()
      })
    })
  })

  xdescribe('create', () => {
    let file = ''

    beforeAll(() => {
      file = this.file
    })

    it('e2e: creates a file', (done) => {
      FileService.create({ file }).then(({ data }) => {
        expect(Object.keys(data)).toEqual([ 'id', 'name' ])
        done()
      }).catch(({ message }) => {
        fail()
        done()
      })
    })
  })
})
