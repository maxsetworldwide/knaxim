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
import Vuex from 'vuex'

import auth from './auth.module'
import file from './file.module'
// import files from './files.module'
import search from './search.module'
import acronyms from './acronyms.module'
import recents from './recents.module'
import folder from './folder.module'
import owner from './owner.module'
import group from './group.module'
import preview from './preview.module'
import error from './error.module'
import nlp from './nlp.module'
import evnt from './event.module'

import UserService from '@/service/user'
import { LOAD_SERVER, HANDLE_SERVER_STATE, AFTER_LOGIN } from './actions.type'
import { PROCESS_SERVER_STATE, PUSH_ERROR } from './mutations.type'

Vue.use(Vuex)

// TODO: Much of this root level code should be in or is already apart of
//  the search.module; move all search related code to the search.module.
export default new Vuex.Store({
  strict: true,
  modules: {
    auth,
    file,
    search,
    acronyms,
    recents,
    folder,
    owner,
    group,
    preview,
    error,
    nlp,
    evnt
  },
  // TODO: Extract all the search functionality into a module!
  state: {
    appSideType: null,
    serverLoading: 0
  },

  getters: {
    appSideType (state) {
      return state.appSideType
    },
    loading (s, g) {
      return (
        g.fileLoading ||
        g.folderLoading ||
        g.authLoading ||
        g.groupLoading ||
        g.ownerLoading ||
        g.searchLoading ||
        s.serverLoading > 0
      )
    }
  },

  actions: {
    async [LOAD_SERVER] ({ commit, dispatch }) {
      try {
        commit('serverloadingchange', 1)
        let res = await UserService.completeProfile()
        dispatch(HANDLE_SERVER_STATE, res.data)
        commit(PROCESS_SERVER_STATE, res.data)
      } catch (err) {
        commit(PUSH_ERROR, new Error(`LOAD_SERVER: ${err}`))
      } finally {
        commit('serverloadingchange', -1)
      }
    },
    [AFTER_LOGIN] ({ dispatch }) {
      dispatch(LOAD_SERVER)
    }
  },

  mutations: {
    serverloadingchange (state, delta) {
      state.serverLoading += delta
    }
  }
})
