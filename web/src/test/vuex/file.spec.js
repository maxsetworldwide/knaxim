import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import {
  LOAD_SERVER,
  CREATE_FILE,
  DELETE_FILES,
  GET_FILE,
  CREATE_WEB_FILE
} from '@/store/actions.type'
import {
  FILE_LOADING,
  SET_FILE,
  PROCESS_SERVER_STATE
} from '@/store/mutations.type'
import modul from '@/store/file.module'

describe('File store', function () {
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
    it('adjusts file loading', function () {
      m[FILE_LOADING](this.state, 5)
      expect(this.state.loading).toBe(5)
      m[FILE_LOADING](this.state, -5)
      expect(this.state.loading).toBe(0)
    })
    it('saves file by id', function () {
      m[SET_FILE](this.state, {
        id: 'id',
        data: 'data'
      })
      expect(this.state.fileSet.id.data).toBe('data')
    })
    it('process server state', function () {
      m[PROCESS_SERVER_STATE](this.state, {
        files: {
          id: { id: 'id', data: 'data' }
        },
        public: [ 'public' ],
        user: { files: { own: [ 'own' ], view: [ 'view' ] } },
        groups: {
          first: {
            files: {
              own: [ '1own' ],
              view: [ '1view' ]
            }
          },
          second: {
            files: {
              own: [ '2own' ],
              view: [ '2view' ]
            }
          }
        }
      })
      expect(this.state.fileSet).toEqual({
        id: { id: 'id', data: 'data' }
      })
      expect(this.state.public).toEqual([ 'public' ])
      expect(this.state.user).toEqual({
        owned: [ 'own' ],
        shared: [ 'view' ]
      })
      expect(this.state.groups).toEqual({
        first: {
          owned: [ '1own' ],
          shared: [ '1view' ]
        },
        second: {
          owned: [ '2own' ],
          shared: [ '2view' ]
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
    it('Gets a File', async function () {
      mock.onGet('file/fid').reply(200, {
        file: {
          id: 'fid',
          data: 'data'
        },
        size: 11,
        count: 12
      })
      await testAction(
        a[GET_FILE],
        {
          id: 'fid'
        },
        {
          mutations: [
            { type: FILE_LOADING, payload: 1 },
            {
              type: SET_FILE,
              payload: {
                id: 'fid',
                data: 'data',
                size: 11,
                count: 12
              }
            },
            { type: FILE_LOADING, payload: -1 }
          ]
        },
        this.state
      )
      this.state.fileSet.fid = {
        id: 'fid',
        data: 'data',
        size: 11,
        count: 12
      }
      await testAction(
        a[GET_FILE],
        { id: 'fid' },
        {
          resolve: {
            id: 'fid',
            data: 'data',
            size: 11,
            count: 12
          }
        },
        this.state
      )
    })
    it('creates a file', async function () {
      mock.onPut('file', {
        group: 'group',
        dir: 'dir',
        file: 'file'
      }).reply(200, 'success')
      await testAction(
        a[CREATE_FILE],
        {
          file: 'file',
          folder: 'dir',
          group: 'group'
        }, {
          mutations: [
            { type: FILE_LOADING, payload: 1 },
            { type: FILE_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOAD_SERVER }
          ],
          resolve: 'success'
        }
      )
    })
    it('creates a web file', async function () {
      mock.onPut('file/webpage', {
        group: 'group',
        url: 'url',
        dir: 'dir'
      }).reply(200, 'success')
      await testAction(
        a[CREATE_WEB_FILE],
        {
          url: 'url',
          group: 'group',
          folder: 'dir'
        },
        {
          mutations: [
            { type: FILE_LOADING, payload: 1 },
            { type: FILE_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOAD_SERVER }
          ],
          resolve: 'success'
        }
      )
    })
    it('deletes file', async function () {
      mock.onDelete('file/fid').reply(200)
      await testAction(
        a[DELETE_FILES],
        {
          ids: [ 'fid' ]
        },
        {
          mutations: [
            { type: FILE_LOADING, payload: 1 },
            { type: FILE_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOAD_SERVER }
          ]
        }
      )
    })
  })
})
