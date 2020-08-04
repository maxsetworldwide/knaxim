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
import Vue from 'vue'
import FileService from '@/service/file'
import SearchService from '@/service/search'
import {
  SEARCH,
  SEARCH_TAG,
  LOAD_MATCHED_LINES,
  LOAD_FILE_MATCH_LINES
} from './actions.type'
import {
  SEARCH_LOADING,
  PUSH_ERROR,
  NEW_SEARCH,
  SET_MATCHED_LINES,
  LOADING_MATCHED_LINES,
  SET_MATCHES,
  DEACTIVATE_SEARCH
} from './mutations.type'

const state = {
  matches: [],
  loading: 0,
  history: [],
  activeSearch: false,
  lines: {},
  summaryStep: 100,
  cancelSearch: false
}

const actions = {
  // TODO: isLoading is updated correctly but not Utalized yet.  If a user spams
  // the FILES_SEARCH action they might get duplicate results.
  // TODO: Use isLoading: true to cancel the current search and start a new one.
  async [SEARCH] ({ commit, dispatch, getters }, { find, acr }) {
    commit(SEARCH_LOADING, 1)
    if (acr) {
      find = `"${find}" ${acr}`
    }
    if (find.length < 1) {
      return false
    }
    try {
      let fileList = await new Promise((resolve, reject) => {
        commit('cancelSearch', () => {
          reject(new Error('search canceled'))
        })
        commit(NEW_SEARCH, { find })
        let context = SearchService.newOwnerContext((getters.activeGroup && getters.activeGroup.id) || getters.currentUser.id)
        let match = SearchService.newMatchCondition(find)
        SearchService.FileTags({ context, match }).then(({ data }) => {
          if (data.matched && data.matched.length > 0) {
            return data.matched.map(item => {
              return {
                ...item.file,
                count: item.count
              }
            })
          }
          return []
        }).then(r => resolve(r)).catch(e => reject(e))
      })
      commit(SET_MATCHES, fileList)
      await dispatch(LOAD_MATCHED_LINES, { find, files: fileList })
    } catch (err) {
      commit(PUSH_ERROR, err.addDebug('action SEARCH'))
    } finally {
      commit(SEARCH_LOADING, -1)
    }
  },
  async [SEARCH_TAG] ({ commit, dispatch }, { context, match }) {
    commit(SEARCH_LOADING, 1)
    try {
      let fileList = await new Promise((resolve, reject) => {
        commit('cancelSearch', () => {
          reject(new Error('search canceled'))
        })
        commit(NEW_SEARCH, { find: match.word })
        SearchService.FileTags({ context, match }).then(({ data }) => {
          if (data.matched && data.matched.length > 0) {
            return data.matched.map(item => {
              return {
                ...item.file,
                count: item.count
              }
            })
          }
        }).then(r => resolve(r)).catch(e => reject(e))
      })
      commit(SET_MATCHES, fileList)
      await dispatch(LOAD_MATCHED_LINES, { find: match.word, files: fileList })
    } catch (err) {
      commit(PUSH_ERROR, err.addDebug('action SEARCH'))
    } finally {
      commit(SEARCH_LOADING, -1)
    }
  },
  async [LOAD_MATCHED_LINES] ({ commit, dispatch }, { find, files }) {
    files.forEach(({ id, count }) => {
      commit(LOADING_MATCHED_LINES, { id, delta: 1 })
      commit(SET_MATCHED_LINES, { id, matched: [] })
      dispatch(LOAD_FILE_MATCH_LINES, { find, id, limit: count })
        .finally(() => commit(LOADING_MATCHED_LINES, { id, delta: -1 }))
    })
  },
  async [LOAD_FILE_MATCH_LINES] ({ commit, state }, { find, id, limit }) {
    try {
      commit(LOADING_MATCHED_LINES, { id, delta: 1 })
      commit(SET_MATCHED_LINES, { id, matched: [] })
      let found = []
      for (let start = 0; start < limit && found.length < 4; start += state.summaryStep) {
        let lines = await FileService.search({
          fid: id,
          start,
          end: start + state.summaryStep,
          find
        }).then(({ data }) => (data.lines || []))
        found = [ ...found, ...lines ]
      }
      commit(SET_MATCHED_LINES, { id, matched: found })
    } catch (err) {
      commit(PUSH_ERROR, new Error(`LOAD_FILE_MATCH_LINES ${err}`))
    } finally {
      commit(LOADING_MATCHED_LINES, { id, delta: -1 })
    }
  }
}

const mutations = {
  [SEARCH_LOADING] (state, delta) {
    state.loading += delta
  },
  [NEW_SEARCH] (state, { find }) {
    state.history = state.history.filter(h => h !== find)
    if (state.history.unshift(find) > 10) {
      state.history.pop()
    }
    state.activeSearch = true
    state.matches = []
  },
  [DEACTIVATE_SEARCH] (state) {
    state.activeSearch = false
  },
  [SET_MATCHES] (state, matches) {
    state.matches = matches
  },
  [LOADING_MATCHED_LINES] (state, { id, delta }) {
    if (!state.lines[id]) {
      Vue.set(state.lines, id, {
        loadingCount: 0,
        get loading () { return this.loadingCount > 0 },
        matched: []
      })
    }
    Vue.set(state.lines[id], 'loadingCount', state.lines[id].loadingCount + delta)
  },
  [SET_MATCHED_LINES] (state, { id, matched }) {
    if (!state.lines[id]) {
      Vue.set(state.lines, id, {
        loadingCount: 0,
        get loading () { return this.loadingCount > 0 },
        matched: []
      })
    }
    Vue.set(state.lines[id], 'matched', matched)
  },
  cancelSearch (state, newfunc) {
    if (state.cancelSearch) {
      state.cancelSearch()
    }
    state.cancelSearch = newfunc
  }
}

const getters = {
  searchMatches (state) {
    return state.activeSearch ? state.matches : []
  },
  currentSearch (state) {
    return state.activeSearch ? state.history[0] || '' : ''
  },
  searchHistory (state) {
    return state.history
  },
  searchLoading (state) {
    return state.loading > 0
  },
  searchLines (state) {
    return state.lines
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
