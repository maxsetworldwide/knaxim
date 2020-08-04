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
    it('Populates Associated Acronyms', async function () {
      let mock = new MockAdapter(axios)
      mock.onGet('/acronym/aaa').reply(200, { matched: ['triple A'] }, { 'Content-Type': 'application/json' })
      await testAction(
        a[ACRONYMS],
        { acronym: 'aaa' },
        {
          mutations: [
            { type: LOADING_ACRONYMS, payload: 1 },
            { type: SET_ACRONYMS, payload: { acronyms: ['triple A'] } },
            { type: LOADING_ACRONYMS, payload: -1 }
          ]
        }
      )
      mock.restore()
    })
    it('Populates as empty with non string value', async function () {
      await testAction(
        a[ACRONYMS],
        { acronym: 5 },
        { mutations: [{ type: SET_ACRONYMS, payload: { acronyms: [] } }] }
      )
    })
  })
})
