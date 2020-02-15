import FileService from '@/service/file'
import { FILE_SLICES, CREATE_FILE } from './actions.type'
import {
  START_SLICES,
  END_SLICES,
  FILE_CREATED,
  FILE_START_LOADING,
  FILE_STOP_LOADING,
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
    context.commit(FILE_START_LOADING)
    return FileService.create(params)
      .then(res => res.data)
      .finally(() => { context.commit(FILE_STOP_LOADING) })
  }
}

const mutations = {
  [FILE_START_LOADING] (state) {
    state.loading += 1
  },
  [FILE_STOP_LOADING] (state) {
    state.loading -= 1
  },
  [SET_FILE] (state, file) {
    state.fileSet[file.id] = file
  },
  [PROCESS_SERVER_STATE] (state, server) {
    state.fileSet = server.files
    state.public = server.public
    state.user.owned = server.user.files.own
    state.user.shared = server.user.files.view
    state.groups = {}
    for (let key in server.groups) {
      state.groups[key] = {
        owned: server.groups[key].files.own,
        shared: server.groups[key].files.view
      }
    }
  }
}

const getters = {
  fileLoading (state) {
    return state.loading > 0
  },
  populateFiles (state) {
    return id => {
      if (typeof id === 'string') {
        return state.fileSet[id]
      }
      if (id instanceof Array) {
        return id.map(i => {
          return state.fileSet[i]
        })
      }
    }
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
