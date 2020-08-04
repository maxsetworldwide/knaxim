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
import FileService from '@/service/file'
import Vue from 'vue'
import { LOAD_PREVIEW } from './actions.type'
import { LOADING_PREVIEW, SET_PREVIEW, PUSH_ERROR } from './mutations.type'

const state = {
  preview: {}
}

const getters = {
  filePreview (state) {
    return state.preview
  }
}

const actions = {
  async [LOAD_PREVIEW] ({ commit, getters }, { id }) {
    commit(LOADING_PREVIEW, { id, delta: 1 })
    let lines = []
    try {
      if (!getters.filePreview[id].lines) {
        let data = await FileService.slice({
          fid: id,
          start: 0,
          end: 3
        }).then(res => res.data || {}).then(d => d.lines || [])
        lines = data.map(d => {
          return d.Content[0]
        })
        commit(SET_PREVIEW, { id, lines })
      } else {
        lines = getters.filePreview[id].lines
      }
    } catch (err) {
      commit(SET_PREVIEW, { id, lines: ['Unable to load preview of file.', err] })
      commit(PUSH_ERROR, err.addDebug('action LOAD_PREVIEW'))
    } finally {
      commit(LOADING_PREVIEW, { id, delta: -1 })
    }
    return lines
  }
}

const mutations = {
  [LOADING_PREVIEW] (state, { id, delta }) {
    if (!state.preview[id]) {
      Vue.set(state.preview, id, {
        loadingCount: 0,
        get loading () { return this.loadingCount > 0 },
        lines: null
      })
    }
    Vue.set(state.preview[id], 'loadingCount', state.preview[id].loadingCount + delta)
  },
  [SET_PREVIEW] (state, { id, lines }) {
    if (!state.preview[id]) {
      Vue.set(state.preview, id, {
        loadingCount: 0,
        get loading () { return this.loadingCount > 0 },
        lines: null
      })
    }
    Vue.set(state.preview[id], 'lines', lines)
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
