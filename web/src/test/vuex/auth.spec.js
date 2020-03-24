import modul from '@/store/auth.module'
import {
  LOGIN,
  AFTER_LOGIN,
  LOGOUT,
  REGISTER,
  GET_USER,
  CHANGE_PASSWORD,
  SEND_RESET_REQUEST,
  RESET_PASSWORD
} from '@/store/actions.type'
import {
  SET_USER,
  PURGE_AUTH,
  PROCESS_SERVER_STATE,
  AUTH_LOADING
} from '@/store/mutations.type'
import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'

describe('Authentication store', function () {
  describe('Mutations', function () {
    const m = modul.mutations
    beforeEach(function () {
      this.state = {
        user: null,
        loading: 0
      }
    })
    it('assigns a user', function () {
      m[SET_USER](this.state, 'user')
      expect(this.state.user).toBe('user')
    })
    it('deletes user', function () {
      this.state.user = 'delete me'
      m[PURGE_AUTH](this.state)
      expect(this.state.user).toBeFalsy()
    })
    it('processes server state', function () {
      m[PROCESS_SERVER_STATE](this.state, { user: {
        id: 'id',
        name: 'name',
        data: 'data'
      } })
    })
    it('tracks loading state', function () {
      m[AUTH_LOADING](this.state, 1)
      expect(this.state.loading).toBe(1)
      m[AUTH_LOADING](this.state, -1)
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
    it('logs in', async function () {
      mock.onPost('user/login', {
        name: 'testuser',
        pass: 'passtest'
      })
        .reply(200,
          {
            id: 'id'
          }
        )
      await testAction(
        a[LOGIN],
        {
          login: 'testuser',
          password: 'passtest'
        },
        {
          mutations: [
            { type: PURGE_AUTH },
            { type: AUTH_LOADING, payload: 1 },
            { type: AUTH_LOADING, payload: -1 }
          ],
          actions: [
            { type: AFTER_LOGIN }
          ],
          resolve: {
            id: 'id'
          }
        }
      )
    })
    it('logs out', async function () {
      mock.onDelete('user').reply(200)
      await testAction(
        a[LOGOUT],
        {},
        {
          mutations: [
            { type: PURGE_AUTH },
            { type: AUTH_LOADING, payload: 1 },
            { type: AUTH_LOADING, payload: -1 }
          ]
        }
      )
    })
    it('creates new user', async function () {
      mock.onPut('user', {
        name: 'name',
        pass: 'pass',
        email: 'email'
      }).reply(
        200,
        { id: 'id' }
      )
      await testAction(
        a[REGISTER],
        {
          email: 'email',
          login: 'name',
          password: 'pass'
        },
        {
          mutations: [
            { type: AUTH_LOADING, payload: 1 },
            { type: AUTH_LOADING, payload: -1 }
          ],
          resolve: { id: 'id' }
        }
      )
    })
    it('changes password', async function () {
      mock.onPost('user/pass', {
        oldpass: 'oldpass',
        newpass: 'newpass'
      }).reply(200)
      await testAction(
        a[CHANGE_PASSWORD],
        {
          oldpass: 'oldpass',
          newpass: 'newpass'
        },
        {
          mutations: [
            { type: AUTH_LOADING, payload: 1 },
            { type: AUTH_LOADING, payload: -1 }
          ],
          actions: [
            { type: LOGOUT }
          ]
        }
      )
    })
    it('sends reset request', async function () {
      mock.onPut('user/reset', {
        name: 'name'
      }).reply(200)
      await testAction(
        a[SEND_RESET_REQUEST],
        {
          name: 'name'
        },
        {
          mutations: [
            { type: AUTH_LOADING, payload: 1 },
            { type: AUTH_LOADING, payload: -1 }
          ]
        }
      )
    })
    it('reset password', async function () {
      mock.onPost('user/reset', {
        key: 'key',
        newpass: 'newpass'
      }).reply(200)
      await testAction(
        a[RESET_PASSWORD],
        {
          passkey: 'key',
          newpass: 'newpass'
        },
        {
          mutations: [
            { type: AUTH_LOADING, payload: 1 },
            { type: AUTH_LOADING, payload: -1 }
          ]
        }
      )
    })
    it('get user', async function () {
      mock.onGet('user').reply(200, { id: 'id' })
      await testAction(
        a[GET_USER],
        {
          quiet: true
        },
        {
          mutations: [
            { type: AUTH_LOADING, payload: 1 },
            { type: SET_USER, payload: { id: 'id' } },
            { type: AUTH_LOADING, payload: -1 }
          ],
          resolve: { id: 'id' }
        }
      )
    })
  })
})
