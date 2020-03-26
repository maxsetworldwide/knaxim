import { TOUCH } from '@/store/mutations.type'
import modul from '@/store/recents.module'

describe('Recent Module', function () {
  beforeAll(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  beforeEach(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  describe('Mutations', function () {
    const m = modul.mutations
    it('Records files', function () {
      m[TOUCH](this.state, 'fid')
      expect(this.state.files).toEqual(['fid'])
    })
    it('touched files are moved to the front', function () {
      this.state.files = ['1', '2']
      m[TOUCH](this.state, '2')
      expect(this.state.files).toEqual(['2', '1'])
    })
  })
})
