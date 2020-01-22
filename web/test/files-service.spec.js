import SearchService from '@/service/search'
import FileService from '@/service/file'
import setDefaults from '../test/setup'

describe('FilesService', function () {
  // Test Setup
  beforeAll(async () => {
    await setDefaults(this)
  })

  /* Begin Tests */
  describe('search', () => {
    let find = ''
    beforeAll(() => {
      find = this.find
    })

    it('e2e: returns [matched] property', (done) => {
      SearchService.userFiles({ find })
        .then(({ data }) => {
          expect(Object.keys(data)).toContain(jasmine.arrayContaining([
            'matched'
          ]))
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
