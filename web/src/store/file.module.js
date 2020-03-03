import FileService from '@/service/file'
import {
  LOAD_SERVER,
  CREATE_FILE,
  DELETE_FILES,
  GET_FILE,
  CREATE_WEB_FILE
} from './actions.type'
import {
  FILE_LOADING,
  SET_FILE,
  PROCESS_SERVER_STATE,
  PUSH_ERROR
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
  async [GET_FILE] ({ commit, state }, { id, overwrite = false }) {
    if (overwrite || !state.fileSet[id]) {
      commit(FILE_LOADING, 1)
      try {
        let file = await FileService.info({ fid: id }).then(res => {
          return {
            size: res.data.size || 0,
            count: res.data.count || 0,
            ...res.data.file
          }
        })
        commit(SET_FILE, file)
      } catch (err) {
        commit(PUSH_ERROR, new Error(`GET_FILE: ${err}`))
      } finally {
        commit(FILE_LOADING, -1)
      }
    }
    return state.fileSet[id]
  },
  [CREATE_FILE] (context, params) {
    context.commit(FILE_LOADING, 1)
    return FileService.create(params)
      .then(res => res.data)
      .catch(err => context.commit(PUSH_ERROR, new Error(`CREATE_FILE: ${err}`)))
      .finally(() => context.dispatch(LOAD_SERVER))
      .finally(() => {
        context.commit(FILE_LOADING, -1)
      })
  },
  [CREATE_WEB_FILE] ({ commit, dispatch }, { url, group, folder }) {
    commit(FILE_LOADING, 1)
    return FileService.webpage({ url, group, folder })
      .then(res => res.data)
      .catch(err => commit(PUSH_ERROR, new Error(`CREATE_WEB_FILE: ${err}`)))
      .finally(() => dispatch(LOAD_SERVER))
      .finally(() => commit(FILE_LOADING, -1))
  },
  [DELETE_FILES] ({ commit, dispatch }, { ids }) {
    commit(FILE_LOADING, 1)
    return Promise.allSettled(
      (ids || []).map(
        id => FileService.erase({ fid: id }).catch(err => { commit(PUSH_ERROR, new Error(`DELETE_FILES{${id}}: ${err}`)); return err })
      )
    )
      .finally(() => {
        dispatch(LOAD_SERVER)
      })
      .finally(() => {
        commit(FILE_LOADING, -1)
      })
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
  sharedFiles (state, getters) {
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
