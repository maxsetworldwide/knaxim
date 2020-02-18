import Vue from 'vue'
import FolderService from '@/service/folder'
import { LOAD_FOLDERS, LOAD_FOLDER, PUT_FILE_FOLDER, REMOVE_FILE_FOLDER, HANDLE_SERVER_STATE } from './actions.type'
import { FOLDER_LOADING, SET_FOLDER, FOLDER_ADD, FOLDER_REMOVE, ACTIVATE_GROUP, ACTIVATE_FOLDER, DEACTIVATE_FOLDER } from './mutations.type'

const state = {
  user: {}, // map filename to list of fileids
  loading: 0, // when greater then 0 folders are being loaded
  group: {}, // map[group id]map[filename][]fileid
  active: []
}

const actions = {
  async [LOAD_FOLDERS] (context, { group, overwrite }) {
    context.commit(FOLDER_LOADING, 1)
    try {
      var names = await FolderService.list({ group }).then(({ data }) => {
        return data.folders
      })
      if (!names) {
        names = []
      }

      await Promise.all(names.map((name) => {
        return context.dispatch(LOAD_FOLDER, { name, group, overwrite })
      }))
    } catch {
      // TODO: handler err
    }
    context.commit(FOLDER_LOADING, -1)
  },
  async [LOAD_FOLDER] (context, { name, group, overwrite }) {
    if (overwrite || context.getters.getFolder({ name, group }).length < 1) {
      context.commit(FOLDER_LOADING, 1)
      try {
        var response = await FolderService.info({ name, group })
        name = response.data.name || name
        var files = response.data.files || []
        context.commit(SET_FOLDER, {
          group,
          name,
          files
        })
      } catch {
        // TODO: handle Error
      }
      context.commit(FOLDER_LOADING, -1)
    }
  },
  async [PUT_FILE_FOLDER] (context, { fid, name, group }) {
    context.commit(FOLDER_LOADING, 1)
    FolderService.add({ fid, name, group }).then(async () => {
      await context.dispatch(LOAD_FOLDER, { group, name, overwrite: true })
    }).finally(() => {
      context.commit(FOLDER_LOADING, -1)
    })
  },
  async [REMOVE_FILE_FOLDER] (context, { fid, name, group }) {
    context.commit(FOLDER_LOADING, 1)
    FolderService.remove({ fid, name, group }).then(async () => {
      await context.dispatch(LOAD_FOLDER, { group, name, overwrite: true })
    }).finally(() => {
      context.commit(FOLDER_LOADING, -1)
    })
  },
  async [HANDLE_SERVER_STATE] ({ commit, dispatch }, { user, groups }) {
    commit(FOLDER_LOADING, 1)
    let proms = (user.folders || []).map(name => dispatch(LOAD_FOLDER, { name }))
    for (let gid in groups) {
      proms.push(...(groups[gid].folders || []).map(name => dispatch(LOAD_FOLDER, { name, group: gid })))
    }
    try {
      await Promise.all(proms)
    } catch {
      // TODO: handle Error
    }
    commit(FOLDER_LOADING, -1)
  }
}

const mutations = {
  [ACTIVATE_GROUP] (context) {
    context.active = []
  },
  [ACTIVATE_FOLDER] (context, name) {
    let newactive = context.active.filter(val => {
      return val !== name
    })
    newactive.unshift(name)
    context.active = newactive
  },
  [DEACTIVATE_FOLDER] (context, name) {
    context.active = context.active.filter(val => val !== name)
  },
  [FOLDER_LOADING] (context, delta) {
    context.loading += delta
  },
  [SET_FOLDER] (context, { group, name, files }) {
    if (!group) {
      Vue.set(context.user, name, files)
      return
    }
    if (!context.group[group]) {
      context.group[group] = {}
    }
    Vue.set(context.group[group], name, files)
  },
  [FOLDER_ADD] (context, { group, name, fid }) {
    if (!group) {
      if (!context.user[name]) {
        context.user[name] = []
      }
      context.user[name].push(fid)
    }
    if (!context.group[group]) {
      context.group[group] = {}
    }
    if (!context.group[group][name]) {
      context.group[group][name] = []
    }
    context.group[group][name].push(fid)
  },
  [FOLDER_REMOVE] (context, { group, name, fid }) {
    if (!group) {
      if (!context.user[name]) {
        return
      }
      let nuser = context.user
      nuser[name] = context.user[name].filter((ele) => {
        return ele !== fid
      })
      context.user = nuser
    }
    if (context.group[group] && context.group[group][name]) {
      let ngroup = context.group
      ngroup[group][name] = context.group[group][name].filter((ele) => {
        return ele !== fid
      })
      context.group = ngroup
    }
  }
}

const getters = {
  folders (state, getters) {
    if (!getters.activeGroup) {
      return state.user
    }
    return state.group[getters.activeGroup]
  },
  getFolder (state, getters) {
    return (name) => {
      if (!getters.activeGroup) {
        return state.user[name] || []
      }
      if (!state.group[getters.activeGroup]) {
        return []
      }
      return state.group[getters.activeGroup][name] || []
    }
  },
  activeFolders ({ active }) {
    return active
  },
  folderLoading ({ loading }) {
    return loading > 0
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
