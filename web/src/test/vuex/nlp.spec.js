import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import modul from '@/store/nlp.module'
import { NLP_DATA } from '@/store/actions.type'
import {
  LOADING_NLP,
  SET_NLP
} from '@/store/mutations.type'

describe('NLP Module', function () {
  beforeAll(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  beforeEach(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  describe('Mutations', function () {
    const m = modul.mutations
    it('records loading state', function () {
      m[LOADING_NLP](this.state, 5)
      expect(this.state.loading).toBe(5)
      m[LOADING_NLP](this.state, -5)
      expect(this.state.loading).toBe(0)
    })
    it('caches nlp data', function () {
      m[SET_NLP](this.state, {
        fid: 'fid',
        start: 2,
        info: ['a', 'b', 'c'],
        category: 't'
      })
      expect(this.state.topics['fid']).toEqual([null, null, 'a', 'b', 'c'])
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
    it('loads nlp data', async function () {
      mock.onGet('nlp/file/fid/t/0/3').reply(200, {
        file: 'fid',
        info: ['a', 'b', 'c']
      })
      await testAction(
        a[NLP_DATA],
        {
          fid: 'fid',
          category: 't',
          start: 0,
          end: 3
        },
        {
          mutations: [
            { type: LOADING_NLP, payload: 1 },
            {
              type: SET_NLP,
              payload: {
                fid: 'fid',
                start: 0,
                info: ['a', 'b', 'c'],
                category: 't'
              }
            },
            { type: LOADING_NLP, payload: -1 }
          ],
          resolve: ['a', 'b', 'c']
        },
        this.state
      )
    })
  })
})
