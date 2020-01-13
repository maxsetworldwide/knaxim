import Vue from 'vue'
import GroupService from '@/service/group'
import { AFTER_LOGIN, REFRESH_GROUPS, CREATE_GROUP } from './actions.type'
import {
  SET_GROUP,
  ACTIVATE_GROUP
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
  availableGroups (state) {
    let teams = []
    for (const id in state.options) {
      teams.push({
        id,
        name: state.options[id]
      })
    }
    return teams
  }
}

export default {
  actions,
  mutations,
  getters,
  state
}
