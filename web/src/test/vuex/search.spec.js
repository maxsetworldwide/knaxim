import { SEARCH, SEARCH_TAG, LOAD_MATCHED_LINES, LOAD_FILE_MATCH_LINES } from '@/store/actions.type'
import {
  SEARCH_LOADING,
  NEW_SEARCH,
  SET_MATCHED_LINES,
  LOADING_MATCHED_LINES,
  SET_MATCHES,
  DEACTIVATE_SEARCH
} from '@/store/mutations.type'
import modul from '@/store/search.module'
import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'

describe('Search module', function () {
  beforeAll(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  beforeEach(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  describe('Mutations', function () {
    const m = modul.mutations
    it('sets loading state', function () {
      m[SEARCH_LOADING](this.state, 5)
      expect(this.state.loading).toBe(5)
      m[SEARCH_LOADING](this.state, -5)
      expect(this.state.loading).toBe(0)
    })
    it('sets a new search', function () {
      m[NEW_SEARCH](this.state, { find: 'find' })
      expect(this.state.history).toEqual(['find'])
      expect(this.state.activeSearch).toBeTrue()
      expect(this.state.matches).toEqual([])
    })
    it('Deactivates search', function () {
      this.state.activeSearch = true
      m[DEACTIVATE_SEARCH](this.state)
      expect(this.state.activeSearch).toBeFalse()
    })
    it('saves matches', function () {
      m[SET_MATCHES](this.state, 'matches')
      expect(this.state.matches).toBe('matches')
    })
    it('tracks loading state for matched lines', function () {
      m[LOADING_MATCHED_LINES](this.state, {
        id: 'id',
        delta: 5
      })
      expect(this.state.lines.id.loadingCount).toBe(5)
      m[LOADING_MATCHED_LINES](this.state, {
        id: 'id',
        delta: -5
      })
      expect(this.state.lines.id.loadingCount).toBe(0)
    })
    it('save matches results', function () {
      m[SET_MATCHED_LINES](this.state, {
        id: 'id',
        matched: ['hi']
      })
      expect(this.state.lines.id.matched).toEqual(['hi'])
    })
    it('cancels search', function () {
      let track = false
      function first () {
        expect(track).toBeFalse()
        track = true
      }
      function second () {
        expect(track).toBeTrue()
        track = false
      }
      m.cancelSearch(this.state, first)
      m.cancelSearch(this.state, second)
      m.cancelSearch(this.state, null)
      expect(track).toBeFalse()
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
    it('searches', async function () {
      mock
        .onPost('search/tags', {
          context: {
            type: 'owner',
            id: 'user'
          },
          match: 'find',
          _sendJSON: true
        }).reply(200, {
          matched: [
            { file: { id: 'fid' }, count: 1 }
          ]
        })
        .onPost('search/tags', {
          context: {
            id: 'gid',
            type: 'owner'
          },
          match: 'find',
          _sendJSON: true
        }).reply(200, {
          matched: [
            { file: { id: 'fid' }, count: 2 }
          ]
        })
      await testAction(
        a[SEARCH],
        {
          find: 'find'
        },
        {
          mutations: [
            { type: SEARCH_LOADING, payload: 1 },
            { type: 'cancelSearch', ignorePayload: true },
            { type: NEW_SEARCH, payload: { find: 'find' } },
            {
              type: SET_MATCHES,
              payload: [
                { id: 'fid', count: 1 }
              ]
            },
            { type: SEARCH_LOADING, payload: -1 }
          ],
          actions: [
            {
              type: LOAD_MATCHED_LINES,
              payload: {
                find: 'find',
                files: [
                  { id: 'fid', count: 1 }
                ]
              }
            }
          ]
        },
        {},
        {
          activeGroup: null,
          currentUser: {
            id: 'user'
          }
        }
      )
      await testAction(
        a[SEARCH],
        {
          find: 'find'
        },
        {
          mutations: [
            { type: SEARCH_LOADING, payload: 1 },
            { type: 'cancelSearch', ignorePayload: true },
            { type: NEW_SEARCH, payload: { find: 'find' } },
            {
              type: SET_MATCHES,
              payload: [
                { id: 'fid', count: 2 }
              ]
            },
            { type: SEARCH_LOADING, payload: -1 }
          ],
          actions: [
            {
              type: LOAD_MATCHED_LINES,
              payload: {
                find: 'find',
                files: [
                  { id: 'fid', count: 2 }
                ]
              }
            }
          ]
        },
        {},
        {
          activeGroup: { id: 'gid' }
        }
      )
    })
    it('searches tags', async function () {
      mock
        .onPost('search/tags', {
          _sendJSON: true,
          context: {
            type: 'file',
            id: 'fid'
          },
          match: {
            tagtype: 'test',
            word: 'find'
          }
        }).reply(200, { matched: [
          { file: { id: 'fid' }, count: 1 }
        ] })
      await testAction(
        a[SEARCH_TAG],
        {
          context: {
            type: 'file',
            id: 'fid'
          },
          match: {
            tagtype: 'test',
            word: 'find'
          }
        },
        {
          mutations: [
            { type: SEARCH_LOADING, payload: 1 },
            { type: 'cancelSearch', ignorePayload: true },
            { type: NEW_SEARCH, payload: { find: 'find' } },
            {
              type: SET_MATCHES,
              payload: [
                { id: 'fid', count: 1 }
              ]
            },
            { type: SEARCH_LOADING, payload: -1 }
          ],
          actions: [
            {
              type: LOAD_MATCHED_LINES,
              payload: {
                find: 'find',
                files: [
                  { id: 'fid', count: 1 }
                ]
              }
            }
          ]
        }
      )
    })
    it('loads matching lines', async function () {
      await testAction(
        a[LOAD_MATCHED_LINES],
        {
          find: 'find',
          files: [
            {
              id: 'fid',
              count: 1
            },
            {
              id: 'fid2',
              count: 2
            }
          ]
        },
        {
          mutations: [
            {
              type: LOADING_MATCHED_LINES,
              payload: {
                id: 'fid',
                delta: 1
              }
            },
            {
              type: SET_MATCHED_LINES,
              payload: {
                id: 'fid',
                matched: []
              }
            },
            {
              type: LOADING_MATCHED_LINES,
              payload: {
                id: 'fid2',
                delta: 1
              }
            },
            {
              type: SET_MATCHED_LINES,
              payload: {
                id: 'fid2',
                matched: []
              }
            },
            {
              type: LOADING_MATCHED_LINES,
              payload: {
                id: 'fid',
                delta: -1
              }
            },
            {
              type: LOADING_MATCHED_LINES,
              payload: {
                id: 'fid2',
                delta: -1
              }
            }
          ],
          actions: [
            {
              type: LOAD_FILE_MATCH_LINES,
              payload: {
                find: 'find',
                id: 'fid',
                limit: 1
              }
            },
            {
              type: LOAD_FILE_MATCH_LINES,
              payload: {
                find: 'find',
                id: 'fid2',
                limit: 2
              }
            }
          ]
        }
      )
    })
    it('load files matching lines', async function () {
      mock.onGet('file/fid/search/0/100', {
        params: {
          find: 'find'
        }
      }).reply(200, {
        lines: [
          '1',
          '2',
          '3',
          '4'
        ]
      })
      await testAction(
        a[LOAD_FILE_MATCH_LINES],
        {
          find: 'find',
          id: 'fid',
          limit: 10
        },
        {
          mutations: [
            {
              type: LOADING_MATCHED_LINES,
              payload: {
                id: 'fid',
                delta: 1
              }
            },
            {
              type: SET_MATCHED_LINES,
              payload: {
                id: 'fid',
                matched: []
              }
            },
            {
              type: SET_MATCHED_LINES,
              payload: {
                id: 'fid',
                matched: ['1', '2', '3', '4']
              }
            },
            {
              type: LOADING_MATCHED_LINES,
              payload: {
                id: 'fid',
                delta: -1
              }
            }
          ]
        },
        this.state
      )
    })
  })
})
