import Vue from 'vue'
import GroupService from '@/service/group'
import { AFTER_LOGIN, REFRESH_GROUPS, CREATE_GROUP } from './actions.type'
import {
  SET_GROUP,
  ACTIVATE_GROUP,
  PROCESS_SERVER_STATE
} from './mutations.type'

const state = {
  active: null,
  options: {}
}

const actions = {
  [AFTER_LOGIN] (context) {
    context.dispatch(REFRESH_GROUPS)
  },
  async [REFRESH_GROUPS] (context) {
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
  },
  async [CREATE_GROUP] (context, { name }) {
    await GroupService.create({ name })
    await context.dispatch(REFRESH_GROUPS)
  }
}

const mutations = {
  [SET_GROUP] (state, { id, name }) {
    Vue.set(state.options, id, name)
  },
  [ACTIVATE_GROUP] (state, { id }) {
    state.active = id
  },
  [PROCESS_SERVER_STATE] ({ commit }, { groups }) {
    groups.values().forEach(v => commit(SET_GROUP, v))
  }
}

const getters = {
  activeGroup (state) {
    if (!state.active) return null
    return {
      id: state.active,
      name: state.options[state.active]
    }
  },
  availableGroups ({ options }) {
    return options.keys().map(id => {
      return {
        id,
        name: options[id]
      }
    })
  }
}

export default {
  actions,
  mutations,
  getters,
  state
}
