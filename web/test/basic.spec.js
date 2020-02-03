// import FileService from '@/service/file'
// import setDefaults from '../test/setup'

describe('Simple Test/Test Suite', function () {
  /* Begin Tests */
  describe('test', () => {
    let fid = ''

    beforeAll(() => {
      fid = 'abc123'
    })

    it('fid = abc123', () => {
      expect(fid).toEqual('abc123')
    })
  })
})
