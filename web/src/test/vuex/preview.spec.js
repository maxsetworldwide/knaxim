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

import { LOAD_PREVIEW } from '@/store/actions.type'
import { LOADING_PREVIEW, SET_PREVIEW } from '@/store/mutations.type'
import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import modul from '@/store/preview.module'

describe('Preview Module', function () {
  beforeAll(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  beforeEach(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  describe('Mutations', function () {
    const m = modul.mutations
    it('Preview loading state', function () {
      m[LOADING_PREVIEW](this.state, {
        id: 'id',
        delta: 5
      })
      this.state.preview.id ? expect(this.state.preview.id.loading).toBeTrue() : fail('did not initialize preview record')
      m[LOADING_PREVIEW](this.state, {
        id: 'id',
        delta: -5
      })
      expect(this.state.preview.id.loading).toBeFalse()
    })
    it('Saves a Preview', function () {
      m[SET_PREVIEW](this.state, {
        id: 'id',
        lines: ['a', 'b', 'c']
      })
      this.state.preview.id ? expect(this.state.preview.id.lines).toEqual(['a', 'b', 'c']) : fail('did not initialize preview record')
    })
  })
  describe('Actions', function () {
    const a = modul.actions
    const mock = new MockAdapter(axios)
    afterEach(function () {
      mock.reset()
    })
    afterAll(function () {
      mock.restore()
    })
    it('loads a preview', async function () {
      mock.onGet('file/fid/slice/0/3').reply(200, {
        lines: [
          { Content: ['hello'] },
          { Content: ['world'] },
          { Content: ['im back'] }
        ]
      })
      await testAction(
        a[LOAD_PREVIEW],
        {
          id: 'fid'
        },
        {
          mutations: [
            {
              type: LOADING_PREVIEW,
              payload: {
                id: 'fid',
                delta: 1
              }
            },
            {
              type: SET_PREVIEW,
              payload: {
                id: 'fid',
                lines: ['hello', 'world', 'im back']
              }
            },
            {
              type: LOADING_PREVIEW,
              payload: {
                id: 'fid',
                delta: -1
              }
            }
          ]
        },
        {},
        {
          filePreview: {
            fid: {
              lines: false
            }
          }
        }
      )
    })
  })
})
