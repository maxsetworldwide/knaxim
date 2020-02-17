import Vue from 'vue'
import Vuex from 'vuex'

import auth from './auth.module'
import file from './file.module'
import files from './files.module'
import search from './search.module'
import acronyms from './acronyms.module'
import recents from './recents.module'
import folder from './folder.module'
import owner from './owner.module'
import group from './group.module'

import UserService from '@/service/user'
import { LOAD_SERVER, HANDLE_SERVER_STATE, PROCESS_SERVER_STATE, AFTER_LOGIN } from './actions.type'

Vue.use(Vuex)

// TODO: Much of this root level code should be in or is already apart of
//  the search.module; move all search related code to the search.module.
export default new Vuex.Store({
  strict: true,
  modules: {
    auth,
    file,
    files,
    search,
    acronyms,
    recents,
    folder,
    owner,
    group
  },
  // TODO: Extract all the search functionality into a module!
  state: {
    appSideType: null,
    selected: ''
  },

  getters: {
    getAppSideType (state) {
      return state.appSideType
    }
  },

  mutations: {
    setSearch (state, item) {
      state.selected = item
    }
  },

  actions: {
    async [LOAD_SERVER] ({ commit, dispatch }) {
      let res = await UserService.completeProfile()
      dispatch(HANDLE_SERVER_STATE, res.data)
      commit(PROCESS_SERVER_STATE, res.data)
    },
    [AFTER_LOGIN] ({ dispatch }) {
      dispatch(LOAD_SERVER)
    }
  }
})
