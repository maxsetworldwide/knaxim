import Vue from 'vue'
import GroupService from '@/service/group'
import { AFTER_LOGIN, REFRESH_GROUPS, CREATE_GROUP, LOAD_SERVER, ADD_MEMBER, REMOVE_MEMBER } from './actions.type'
import {
  SET_GROUP,
  ACTIVATE_GROUP,
  PROCESS_SERVER_STATE,
  GROUP_LOADING,
  PUSH_ERROR
} from './mutations.type'

const state = {
  active: null,
  ids: [],
  names: {},
  members: {},
  owners: {},
  loading: 0
}

const actions = {
  [AFTER_LOGIN] (context) {
    context.dispatch(REFRESH_GROUPS)
  },
  async [REFRESH_GROUPS] (context) {
    context.commit(GROUP_LOADING, 1)
    try {
      let data = await GroupService.associated({}).then(res => res.data)
      if (data.own) {
        data.own.forEach(t => {
          context.commit(SET_GROUP, t)
        })
      }
      if (data.member) {
        data.member.forEach(t => {
          context.commit(SET_GROUP, t)
        })
      }
    } catch (err) {
      context.commit(PUSH_ERROR, err.addDebug('action REFRESH_GROUPS'))
    } finally {
      context.commit(GROUP_LOADING, -1)
    }
  },
  async [CREATE_GROUP] (context, { name }) {
    context.commit(GROUP_LOADING, 1)
    try {
      await GroupService.create({ name })
      await context.dispatch(LOAD_SERVER)
    } catch (err) {
      context.commit(PUSH_ERROR, err.addDebug('action CREATE_GROUP'))
    } finally {
      context.commit(GROUP_LOADING, -1)
    }
  },
  async [ADD_MEMBER] ({ commit, dispatch, state }, { gid, newMember }) {
    if (!gid) {
      gid = state.active
    }
    commit(GROUP_LOADING, 1)
    try {
      await GroupService.add({
        gid,
        target: newMember
      })
      dispatch(LOAD_SERVER)
    } catch (err) {
      commit(PUSH_ERROR, err.addDebug('action ADD_MEMBER')
    } finally {
      commit(GROUP_LOADING, -1)
    }
  },
  async [REMOVE_MEMBER] ({ commit, dispatch, state }, { gid, newMember }) {
    if (!gid) {
      gid = state.active
    }
    commit(GROUP_LOADING, 1)
    try {
      await GroupService.remove({
        gid,
        target: newMember
      })
      dispatch(LOAD_SERVER)
    } catch (err) {
      commit(PUSH_ERROR, err.addDebug('action REMOVE_MEMBER'))
    } finally {
      commit(GROUP_LOADING, -1)
    }
  }
}

const mutations = {
  [SET_GROUP] (state, { id, name, owner, members }) {
    if (state.ids.reduce((a, i) => { return a && i !== id }, true)) { state.ids.push(id) }
    Vue.set(state.names, id, name || '')
    Vue.set(state.members, id, members || [])
    Vue.set(state.owners, id, owner || '')
  },
  [ACTIVATE_GROUP] (state, { id }) {
    state.active = id
  },
  [PROCESS_SERVER_STATE] (state, { groups }) {
    state.ids = []
    state.names = {}
    state.members = {}
    state.owners = {}
    for (let id in groups) {
      let { name, owner, members } = groups[id]
      state.ids.push(id)
      Vue.set(state.names, id, name || '')
      Vue.set(state.members, id, members || [])
      Vue.set(state.owners, id, owner || '')
    }
  },
  [GROUP_LOADING] (state, delta) {
    state.loading += delta
  }
}

const getters = {
  activeGroup ({ active, names, members, owners }) {
    if (!active) return null
    return {
      id: active,
      name: names[active],
      members: members[active],
      owner: owners[active]
    }
  },
  availableGroups ({ ids, names, members, owners }) {
    return ids.map(id => {
      return {
        id,
        name: names[id],
        members: members[id],
        owner: owners[id]
      }
    })
  },
  groupLoading ({ loading }) {
    return loading > 0
  }
}

export default {
  actions,
  mutations,
  getters,
  state
}
