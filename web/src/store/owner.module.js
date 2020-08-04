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
import OwnerService from '@/service/owner'
import { LOAD_OWNER, LOOKUP_OWNER } from './actions.type'
import { SET_OWNER_NAME, PROCESS_SERVER_STATE, OWNER_LOADING, PUSH_ERROR } from './mutations.type'

const state = {
  names: {}, // map[ownerid]string
  loading: 0
}

const actions = {
  async [LOAD_OWNER] (context, { id, overwrite }) {
    if (!id) {
      let e = new Error(`LOAD_OWNER: id = ${id}`)
      context.commit(PUSH_ERROR, e)
      throw e
    }
    if (overwrite || !context.state.names[id]) {
      context.commit(OWNER_LOADING, 1)
      context.commit(SET_OWNER_NAME, { id, name: 'loading...' })

      try {
        let name = await OwnerService.id(id).then(r => r.data.name)
        context.commit(SET_OWNER_NAME, { id, name })
        return name
      } catch (e) {
        context.commit(PUSH_ERROR, e.addDebug('action LOAD_OWNER'))
        return 'Unknown'
      } finally {
        context.commit(OWNER_LOADING, -1)
      }
    } else {
      return context.state.names[id]
    }
  },
  async [LOOKUP_OWNER] ({ commit, state }, { name, overwrite = false }) {
    let foundid = null
    if (!overwrite) {
      for (let id in state.names) {
        if (state.names[id] === name) {
          foundid = id
          break
        }
      }
    }
    if (overwrite || !foundid) {
      commit(OWNER_LOADING, 1)
      try {
        foundid = await OwnerService.name(name).then(r => r.data.id)
        commit(SET_OWNER_NAME, {
          id: foundid,
          name
        })
      } catch (e) {
        commit(PUSH_ERROR, e.addDebug('action LOOKUP_OWNER'))
        throw e
      } finally {
        commit(OWNER_LOADING, -1)
      }
    }
    return foundid
  }
}

const mutations = {
  [SET_OWNER_NAME] (state, { id, name }) {
    if (id) {
      Vue.set(state.names, id, name)
    }
  },
  [PROCESS_SERVER_STATE] (state, { user, groups }) {
    if (user.id) {
      Vue.set(state.names, user.id, user.name)
    }
    for (let id in groups) {
      if (id) {
        Vue.set(state.names, id, groups[id].name)
      }
    }
  },
  [OWNER_LOADING] (state, delta) {
    state.loading += delta
  }
}

const getters = {
  ownerNames (state) {
    return state.names
  },
  ownerLoading (state) {
    return state.loading > 0
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
