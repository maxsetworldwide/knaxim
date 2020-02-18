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

import UserService from '@/service/user'
import { LOAD_SERVER, HANDLE_SERVER_STATE, AFTER_LOGIN } from './actions.type'
import { PROCESS_SERVER_STATE } from './mutations.type'

Vue.use(Vuex)

// TODO: Much of this root level code should be in or is already apart of
//  the search.module; move all search related code to the search.module.
export default new Vuex.Store({
  strict: true,
  modules: {
    auth,
    file,
    // files,
    search,
    acronyms,
    recents,
    folder,
    owner,
    group,
    preview
  },
  // TODO: Extract all the search functionality into a module!
  state: {
    appSideType: null
  },

  getters: {
    appSideType (state) {
      return state.appSideType
    },
    loading (s, g) {
      return g.fileLoading ||
      g.folderLoading ||
      g.authLoading ||
      g.groupLoading ||
      g.ownerLoading ||
      g.searchLoading
    }
  },

  mutations: {
    setSearch (state, item) {
      state.selected = item
    }
  },

  actions: {
    async [LOAD_SERVER] ({ commit, dispatch }) {
      try {
        let res = await UserService.completeProfile()
        dispatch(HANDLE_SERVER_STATE, res.data)
        commit(PROCESS_SERVER_STATE, res.data)
      } catch {
        // TODO: handle error
      }
    },
    [AFTER_LOGIN] ({ dispatch }) {
      dispatch(LOAD_SERVER)
    }
  }
})
