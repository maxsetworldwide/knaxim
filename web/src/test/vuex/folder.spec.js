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
import { LOAD_FOLDERS, LOAD_FOLDER, PUT_FILE_FOLDER, REMOVE_FILE_FOLDER, HANDLE_SERVER_STATE, LOAD_SERVER } from '@/store/actions.type'
import { FOLDER_LOADING, SET_FOLDER, FOLDER_ADD, FOLDER_REMOVE, ACTIVATE_GROUP, ACTIVATE_FOLDER, DEACTIVATE_FOLDER } from '@/store/mutations.type'
import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import modul from '@/store/folder.module'

describe('Folder Store', function () {
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
    it('clears active folders on changing groups', function () {
      this.state.active.push('to be removed')
      m[ACTIVATE_GROUP](this.state)
      expect(this.state.active).toEqual([])
    })
    it('activates a folder', function () {
      this.state.active.push('first', 'second')
      m[ACTIVATE_FOLDER](this.state, 'second')
      expect(this.state.active).toEqual(['second', 'first'])
    })
    it('deactivates a folder', function () {
      this.state.active.push('to be removed')
      m[DEACTIVATE_FOLDER](this.state, 'to be removed')
      expect(this.state.active).toEqual([])
    })
    it('adjusts loading state', function () {
      m[FOLDER_LOADING](this.state, 5)
      expect(this.state.loading).toBe(5)
      m[FOLDER_LOADING](this.state, -5)
      expect(this.state.loading).toBe(0)
    })
    it('sets folder', function () {
      m[SET_FOLDER](this.state, {
        name: 'userfolder',
        files: ['1', '2']
      })
      m[SET_FOLDER](this.state, {
        name: 'groupFolder',
        files: ['3', '4'],
        group: 'group'
      })
      expect(this.state.user).toEqual({
        userfolder: ['1', '2']
      })
      expect(this.state.group).toEqual({
        group: {
          groupFolder: ['3', '4']
        }
      })
    })
    it('adds to a folder', function () {
      m[FOLDER_ADD](this.state, {
        name: 'folder',
        fid: 'fid'
      })
      m[FOLDER_ADD](this.state, {
        name: 'folder',
        group: 'group',
        fid: 'fid'
      })
      expect(this.state.user).toEqual({
        folder: ['fid']
      })
      expect(this.state.group).toEqual({
        group: {
          folder: ['fid']
        }
      })
    })
    it('removes from folder', function () {
      this.state.user.folder = ['fid', 'other']
      this.state.group.group = {
        folder: ['fid', 'other']
      }
      m[FOLDER_REMOVE](this.state, {
        name: 'folder',
        fid: 'fid'
      })
      m[FOLDER_REMOVE](this.state, {
        group: 'group',
        name: 'folder',
        fid: 'fid'
      })
      expect(this.state.user).toEqual({
        folder: ['other']
      })
      expect(this.state.group).toEqual({
        group: {
          folder: ['other']
        }
      })
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
    it('gets all folders', async function () {
      mock
        .onGet('dir', { params: { group: 'group' } }).reply(200, {
          folders: ['a', 'b']
        })
        .onGet('dir').reply(200, {
          folders: ['c', 'd']
        })
      await testAction(
        a[LOAD_FOLDERS],
        {},
        {
          mutations: [
            { type: FOLDER_LOADING, payload: 1 },
            { type: FOLDER_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOAD_FOLDER,
              payload: {
                name: 'c',
                group: undefined,
                overwrite: undefined
              }
            },
            { type: LOAD_FOLDER,
              payload: {
                name: 'd',
                group: undefined,
                overwrite: undefined
              }
            }
          ]
        }
      )
      await testAction(
        a[LOAD_FOLDERS],
        { group: 'group', overwrite: true },
        {
          mutations: [
            { type: FOLDER_LOADING, payload: 1 },
            { type: FOLDER_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOAD_FOLDER,
              payload: {
                name: 'a',
                group: 'group',
                overwrite: true
              }
            },
            { type: LOAD_FOLDER,
              payload: {
                name: 'b',
                group: 'group',
                overwrite: true
              }
            }
          ]
        }
      )
    })
    it('loads a folder', async function () {
      mock
        .onGet('dir/a', { params: { group: 'group' } }).reply(200, {
          name: 'a',
          files: ['gfid', 'gfid2']
        })
        .onGet('dir/a').reply(200, {
          name: 'a',
          files: ['fid', 'fid2']
        })
        .onGet('dir/b').reply(200, {
          name: 'b',
          files: ['fid3', 'fid4']
        })
      await testAction(
        a[LOAD_FOLDER],
        {
          name: 'a',
          group: 'group'
        },
        {
          mutations: [
            { type: FOLDER_LOADING, payload: 1 },
            {
              type: SET_FOLDER,
              payload: {
                group: 'group',
                name: 'a',
                files: ['gfid', 'gfid2']
              }
            },
            { type: FOLDER_LOADING, payload: -1 }
          ]
        },
        {},
        {
          getFolder: function () { return [] }
        }
      )
      await testAction(
        a[LOAD_FOLDER],
        {
          name: 'a'
        },
        {
          mutations: [
            { type: FOLDER_LOADING, payload: 1 },
            {
              type: SET_FOLDER,
              payload: {
                group: undefined,
                name: 'a',
                files: ['fid', 'fid2']
              }
            },
            { type: FOLDER_LOADING, payload: -1 }
          ]
        },
        {},
        {
          getFolder: function () { return [] }
        }
      )
      await testAction(
        a[LOAD_FOLDER],
        {
          name: 'b'
        },
        {
          mutations: [
            { type: FOLDER_LOADING, payload: 1 },
            {
              type: SET_FOLDER,
              payload: {
                group: undefined,
                name: 'b',
                files: ['fid3', 'fid4']
              }
            },
            { type: FOLDER_LOADING, payload: -1 }
          ]
        },
        {},
        {
          getFolder: function () { return [] }
        }
      )
    })
    it('adds file to folder', async function () {
      mock
        .onPost('dir/a/content', { group: 'group', id: 'fid' }).reply(200)
      await testAction(
        a[PUT_FILE_FOLDER],
        {
          fid: 'fid',
          name: 'a',
          group: 'group'
        },
        {
          mutations: [
            { type: FOLDER_LOADING, payload: 1 },
            { type: FOLDER_LOADING, payload: -1 }
          ],
          actions: [
            {
              type: LOAD_FOLDER,
              payload: {
                group: 'group',
                name: 'a',
                overwrite: true
              }
            },
            {
              type: LOAD_SERVER
            }
          ]
        }
      )
    })
    it('removes file from folder', async function () {
      mock
        .onDelete('dir/a/content', {
          id: 'fid',
          group: 'group'
        }).reply(200)
      await testAction(
        a[REMOVE_FILE_FOLDER],
        {
          fid: 'fid',
          name: 'a',
          group: 'group'
        },
        {
          mutations: [
            { type: FOLDER_LOADING, payload: 1 },
            { type: FOLDER_LOADING, payload: -1 }
          ],
          actions: [
            {
              type: LOAD_FOLDER,
              payload: {
                group: 'group',
                name: 'a',
                overwrite: true
              }
            },
            { type: LOAD_SERVER }
          ]
        }
      )
    })
    it('handle the server state', async function () {
      await testAction(
        a[HANDLE_SERVER_STATE],
        {
          user: {
            folders: ['a', 'b']
          },
          groups: {
            group: { folders: ['c', 'd'] }
          }
        },
        {
          mutations: [
            { type: FOLDER_LOADING, payload: 1 },
            { type: FOLDER_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOAD_FOLDER, payload: { name: 'a' } },
            { type: LOAD_FOLDER, payload: { name: 'b' } },
            {
              type: LOAD_FOLDER,
              payload: {
                name: 'c',
                group: 'group'
              }
            },
            {
              type: LOAD_FOLDER,
              payload: {
                name: 'd',
                group: 'group'
              }
            }
          ]
        }
      )
    })
  })
})
