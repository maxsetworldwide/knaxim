import modul from '@/store/acronyms.module'
import { SET_ACRONYMS, LOADING_ACRONYMS } from '@/store/mutations.type'
import { ACRONYMS } from '@/store/actions.type'
import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'

describe('Acronym Store', function () {
  describe('Mutations', function () {
    let m = modul.mutations
    it('Adjusts Loading State', function () {
      let state = { loading: 0 }
      m[LOADING_ACRONYMS](state, 5)
      expect(state.loading).toBe(5)
    })

    it('Sets the current acronyms', function () {
      let state = { acronyms: [] }
      m[SET_ACRONYMS](state, { acronyms: ['This is one', '2 I am'] })
      expect(state.acronyms.length).toBe(2)
      expect(state.acronyms[0])
        .toBe('This is one')
      expect(state.acronyms[1]).toBe('2 I am')
    })
  })

  describe('Actions', function () {
    let a = modul.actions
    it('Populates Associated Acronyms', function (done) {
      let mock = new MockAdapter(axios)
      mock.onGet('/acronym/aaa').reply(200, { matched: ['triple A'] }, { 'Content-Type': 'application/json' })
      testAction(
        a[ACRONYMS],
        { acronym: 'aaa' },
        [
          { type: LOADING_ACRONYMS, payload: 1 },
          { type: SET_ACRONYMS, payload: { acronyms: ['triple A'] } },
          { type: LOADING_ACRONYMS, payload: -1 }
        ],
        done
      )
      mock.restore()
    })
  })
})
