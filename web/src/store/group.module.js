import Vue from 'vue'
import GroupService from '@/service/group'
import { AFTER_LOGIN, REFRESH_GROUPS, CREATE_GROUP } from './actions.type'
import {
  SET_GROUP,
  ACTIVATE_GROUP,
  PROCESS_SERVER_STATE,
  GROUP_LOADING
} from './mutations.type'

const state = {
  active: null,
  ids: [],
  names: {},
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
    } catch {
      // TODO: process error
    } finally {
      context.commit(GROUP_LOADING, -1)
    }
  },
  async [CREATE_GROUP] (context, { name }) {
    context.commit(GROUP_LOADING, 1)
    try {
      await GroupService.create({ name })
      await context.dispatch(REFRESH_GROUPS)
    } catch {
      // TODO: handle error
    } finally {
      context.commit(GROUP_LOADING, -1)
    }
  }
}

const mutations = {
  [SET_GROUP] (state, { id, name }) {
    if (state.ids.reduce((a, i) => { return a && i !== id }, true)) { state.ids.push(id) }
    Vue.set(state.names, id, name)
  },
  [ACTIVATE_GROUP] (state, { id }) {
    state.active = id
  },
  [PROCESS_SERVER_STATE] ({ commit }, { groups }) {
    for (let gid in groups) {
      commit(SET_GROUP, groups[gid])
    }
  },
  [GROUP_LOADING] (state, delta) {
    state.loading += delta
  }
}

const getters = {
  activeGroup (state) {
    if (!state.active) return null
    return {
      id: state.active,
      name: state.names[state.active]
    }
  },
  availableGroups ({ ids, names }) {
    return ids.map(id => {
      return {
        id,
        name: names[id]
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
