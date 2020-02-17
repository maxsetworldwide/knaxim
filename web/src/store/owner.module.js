import Vue from 'vue'
import UserService from '@/service/user'
import GroupService from '@/service/group'
import { LOAD_OWNER } from './actions.type'
import { SET_OWNER_NAME, PROCESS_SERVER_STATE, OWNER_LOADING } from './mutations.type'

const state = {
  names: {}, // map[ownerid]string
  loading: 0
}

const actions = {
  async [LOAD_OWNER] (context, { id, overwrite }) {
    if (overwrite || !context.state.names[id]) {
      context.commit(OWNER_LOADING, 1)
      context.commit(SET_OWNER_NAME, { id, name: 'loading...' })
      let response = null
      try {
        response = await UserService.info({ id })
      } catch {
        response = await GroupService.info({ id })
      }
      context.commit(SET_OWNER_NAME, { id, name: response.data.name })
      context.commit(OWNER_LOADING, -1)
      return response.data.name
    } else {
      return context.state.names[id]
    }
  }
}

const mutations = {
  [SET_OWNER_NAME] (state, { id, name }) {
    Vue.set(state.names, id, name)
  },
  [PROCESS_SERVER_STATE] (state, { user, groups }) {
    Vue.set(state.names, user.id, user.name)
    for (let id in groups) {
      Vue.set(state.names, id, groups[id].name)
    }
  },
  [OWNER_LOADING] (state, delta) {
    state.loading += delta
  }
}

const getters = {
  ownerNames (state) {
    return state.names
  },
  ownerLoading (state, delta) {
    state.loading += delta
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
