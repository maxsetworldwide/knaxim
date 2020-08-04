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
import {
  LOAD_OWNER,
  LOOKUP_OWNER
} from '@/store/actions.type'
import {
  SET_OWNER_NAME,
  PROCESS_SERVER_STATE,
  OWNER_LOADING
} from '@/store/mutations.type'
import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import modul from '@/store/owner.module'

describe('Owner Module', function () {
  beforeAll(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  beforeEach(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  describe('Mutations', function () {
    const m = modul.mutations
    it('sets owner name', function () {
      m[SET_OWNER_NAME](this.state, {
        id: 'id',
        name: 'name'
      })
      expect(this.state.names).toEqual({
        id: 'name'
      })
    })
    it('process server state', function () {
      m[PROCESS_SERVER_STATE](this.state, {
        user: {
          id: 'uid',
          name: 'user'
        },
        groups: {
          gid: {
            id: 'gid',
            name: 'group'
          }
        }
      })
      expect(this.state.names).toEqual({
        uid: 'user',
        gid: 'group'
      })
    })
    it('owner loading state', function () {
      m[OWNER_LOADING](this.state, 5)
      expect(this.state.loading).toBe(5)
      m[OWNER_LOADING](this.state, -5)
      expect(this.state.loading).toBe(0)
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
    it('Loads Owner', async function () {
      mock.onGet('owner/id/uid', {}).reply(200, { name: 'user' })
        .onGet('owner/id/gid').reply(200, { name: 'group' })
      await testAction(
        a[LOAD_OWNER],
        {
          id: 'uid'
        },
        {
          mutations: [
            { type: OWNER_LOADING, payload: 1 },
            {
              type: SET_OWNER_NAME,
              payload: {
                id: 'uid',
                name: 'loading...'
              }
            },
            {
              type: SET_OWNER_NAME,
              payload: {
                id: 'uid',
                name: 'user'
              }
            },
            { type: OWNER_LOADING, payload: -1 }
          ],
          resolve: 'user'
        },
        {
          names: {}
        }
      )
      await testAction(
        a[LOAD_OWNER],
        {
          id: 'gid'
        },
        {
          mutations: [
            { type: OWNER_LOADING, payload: 1 },
            {
              type: SET_OWNER_NAME,
              payload: {
                id: 'gid',
                name: 'loading...'
              }
            },
            {
              type: SET_OWNER_NAME,
              payload: {
                id: 'gid',
                name: 'group'
              }
            },
            { type: OWNER_LOADING, payload: -1 }
          ],
          resolve: 'group'
        },
        {
          names: {}
        }
      )
    })
    it('Lookup Owner', async function () {
      mock.onGet('owner/name/user').reply(200, {
        id: 'uid',
        name: 'user'
      })
        .onGet('owner/name/group').reply(200, {
          id: 'gid',
          name: 'group'
        })
      await testAction(
        a[LOOKUP_OWNER],
        {
          name: 'user'
        },
        {
          mutations: [
            { type: OWNER_LOADING, payload: 1 },
            {
              type: SET_OWNER_NAME,
              payload: {
                id: 'uid',
                name: 'user'
              }
            },
            { type: OWNER_LOADING, payload: -1 }
          ],
          resolve: 'uid'
        },
        {
          names: {}
        }
      )
      await testAction(
        a[LOOKUP_OWNER],
        {
          name: 'group'
        },
        {
          mutations: [
            { type: OWNER_LOADING, payload: 1 },
            {
              type: SET_OWNER_NAME,
              payload: {
                id: 'gid',
                name: 'group'
              }
            },
            { type: OWNER_LOADING, payload: -1 }
          ],
          resolve: 'gid'
        }
      )
    })
  })
})
