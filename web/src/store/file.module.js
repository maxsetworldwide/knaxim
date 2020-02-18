import FileService from '@/service/file'
import { CREATE_FILE } from './actions.type'
import {
  FILE_LOADING,
  SET_FILE,
  PROCESS_SERVER_STATE
} from './mutations.type'

const state = {
  loading: 0,
  fileSet: {},
  user: {
    owned: [],
    shared: []
  },
  groups: {},
  public: []
}

const actions = {
  [CREATE_FILE] (context, params) {
    context.commit(FILE_LOADING, 1)
    return FileService.create(params)
      .then(res => res.data)
      .finally(() => { context.commit(FILE_LOADING, -1) })
  }
}

const mutations = {
  [FILE_LOADING] (state, d) {
    state.loading += d
  },
  [SET_FILE] (state, file) {
    state.fileSet[file.id] = file
  },
  [PROCESS_SERVER_STATE] (state, server) {
    state.fileSet = server.files
    state.public = server.public
    state.user.owned = server.user.files.own || []
    state.user.shared = server.user.files.view || []
    state.groups = {}
    for (let key in server.groups) {
      state.groups[key] = {
        owned: server.groups[key].files.own || [],
        shared: server.groups[key].files.view || []
      }
    }
  }
}

const getters = {
  fileLoading (state) {
    return state.loading > 0
  },
  populateFiles (state) {
    return function (id) {
      // console.log('id')
      // console.log(id)
      if (typeof id === 'string') {
        return state.fileSet[id]
      }
      if (id instanceof Array) {
        return id.map(i => {
          return state.fileSet[i]
        })
      }
      throw new Error('unrecognized file id type')
    }
  },
  ownedFiles (state, getters) {
    if (!getters.activeGroup) {
      return state.user.owned
    } else {
      return state.groups[getters.activeGroup.id].owned
    }
  },
  sharedFiles (state) {
    if (!getters.activeGroup) {
      return state.user.shared
    } else {
      return state.groups[getters.activeGroup.id].shared
    }
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
