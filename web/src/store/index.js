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

Vue.use(Vuex)

// TODO: Much of this root level code should be in or is already apart of
//  the search.module; move all search related code to the search.module.
export default new Vuex.Store({
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
  }
})
