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
  AFTER_LOGIN,
  REFRESH_GROUPS,
  CREATE_GROUP,
  LOAD_SERVER,
  ADD_MEMBER,
  REMOVE_MEMBER
} from '@/store/actions.type'
import {
  SET_GROUP,
  ACTIVATE_GROUP,
  PROCESS_SERVER_STATE,
  GROUP_LOADING
} from '@/store/mutations.type'
import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import modul from '@/store/group.module'

describe('Group Module', function () {
  beforeAll(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
    // this.getters = {}
    // for (const getter in modul.getter)
  })
  beforeEach(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  describe('Mutations', function () {
    const m = modul.mutations
    it('saves a group', function () {
      m[SET_GROUP](this.state, {
        id: 'group',
        name: 'name',
        owner: 'owner',
        members: ['devon', 'theresa']
      })
      expect(this.state.ids).toEqual([ 'group' ])
      expect(this.state.names).toEqual({ group: 'name' })
      expect(this.state.members).toEqual({
        group: ['devon', 'theresa']
      })
      expect(this.state.owners).toEqual({
        group: 'owner'
      })
    })
    it('sets an active group', function () {
      m[ACTIVATE_GROUP](this.state, {
        id: 'group'
      })
      expect(this.state.active).toBe('group')
    })
    it('processes server state', function () {
      m[PROCESS_SERVER_STATE](this.state, {
        groups: {
          group: {
            id: 'group',
            name: 'name',
            members: ['devon', 'theresa'],
            owner: 'owner'
          },
          group2: {
            id: 'group2',
            name: 'name2',
            members: ['cat', 'dog'],
            owner: 'piper'
          }
        }
      })
      expect(this.state.ids).toEqual(['group', 'group2'])
      expect(this.state.names).toEqual({
        group: 'name',
        group2: 'name2'
      })
      expect(this.state.owners).toEqual({
        group: 'owner',
        group2: 'piper'
      })
      expect(this.state.members).toEqual({
        group: ['devon', 'theresa'],
        group2: ['cat', 'dog']
      })
    })
    it('group loading state', function () {
      m[GROUP_LOADING](this.state, 5)
      expect(this.state.loading).toBe(5)
      m[GROUP_LOADING](this.state, -5)
      expect(this.state.loading).toBe(0)
    })
  })
  describe('Action', function () {
    const a = modul.actions
    const mock = new MockAdapter(axios)
    afterEach(function () {
      mock.reset()
    })
    afterAll(function () {
      mock.restore()
    })
    it('After Login', async function () {
      await testAction(
        a[AFTER_LOGIN],
        {},
        {
          actions: [
            { type: REFRESH_GROUPS }
          ]
        }
      )
    })
    it('Refreshing groups', async function () {
      mock.onGet('group/options').reply(200, {
        own: ['a', 'b'],
        member: ['c', 'd']
      })
      await testAction(
        a[REFRESH_GROUPS],
        {},
        {
          mutations: [
            { type: GROUP_LOADING, payload: 1 },
            { type: SET_GROUP, payload: 'a' },
            { type: SET_GROUP, payload: 'b' },
            { type: SET_GROUP, payload: 'c' },
            { type: SET_GROUP, payload: 'd' },
            { type: GROUP_LOADING, payload: -1 }
          ]
        }
      )
    })
    it('create a group', async function () {
      mock.onPut('group', {
        newname: 'name'
      }).reply(200)
      await testAction(
        a[CREATE_GROUP],
        { name: 'name' },
        {
          mutations: [
            { type: GROUP_LOADING, payload: 1 },
            { type: GROUP_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOAD_SERVER }
          ]
        }
      )
    })
    it('adds member', async function () {
      mock.onPost('group/grp/member', {
        id: 'new'
      }).reply(200)
      await testAction(
        a[ADD_MEMBER],
        {
          gid: 'grp',
          newMember: 'new'
        },
        {
          mutations: [
            { type: GROUP_LOADING, payload: 1 },
            { type: GROUP_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOAD_SERVER }
          ]
        }
      )
    })
    it('removes member', async function () {
      mock.onDelete('group/grp/member', {
        id: 'new'
      }).reply(200)
      await testAction(
        a[REMOVE_MEMBER],
        {
          gid: 'grp',
          newMember: 'new'
        },
        {
          mutations: [
            { type: GROUP_LOADING, payload: 1 },
            { type: GROUP_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOAD_SERVER }
          ]
        }
      )
    })
  })
})
