import Vue from 'vue'
import OwnerService
import { LOAD_OWNER, LOOKUP_OWNER } from './actions.type'
import { SET_OWNER_NAME, PROCESS_SERVER_STATE, OWNER_LOADING, PUSH_ERROR } from './mutations.type'

const state = {
  names: {}, // map[ownerid]string
  loading: 0
}

const actions = {
  async [LOAD_OWNER] (context, { id, overwrite }) {
    if (!id) {
      let e = new Error(`LOAD_OWNER: id = ${id}`)
      context.commit(PUSH_ERROR, e)
      throw e
    }
    if (overwrite || !context.state.names[id]) {
      context.commit(OWNER_LOADING, 1)
      context.commit(SET_OWNER_NAME, { id, name: 'loading...' })

      try {
        name = await OwnerService.id(id).then(r => r.data.name)
        context.commit(SET_OWNER_NAME, { id, name })
        return name
      } catch (e) {
        commit(PUSH_ERROR, e)
        return 'Unknown'
      } finally {
        commit(OWNER_LOADING, -1)
      }
    } else {
      return context.state.names[id]
    }
  },
  async [LOOKUP_OWNER] ({ commit, state }, { name, overwrite = false }) {
    let foundid = null
    if (!overwrite) {
      for (let id in state.names) {
        if (state.names[id] === name) {
          foundid = id
          break
        }
      }
    }
    if (overwrite || !foundid) {
      commit(OWNER_LOADING, 1)
      try {
        foundid = await OwnerService.name(name).then(r => r.data.id)
        commit(SET_OWNER_NAME, {
          id: foundid,
          name
        })
      } catch (e) {
        commit(PUSH_ERROR, new Error(`LOOKUP_OWNER: ${e}`))
        throw e
      } finally {
        commit(OWNER_LOADING, -1)
      }
    }
    return foundid
  }
}

const mutations = {
  [SET_OWNER_NAME] (state, { id, name }) {
    if (id) {
      Vue.set(state.names, id, name)
    }
  },
  [PROCESS_SERVER_STATE] (state, { user, groups }) {
    if (user.id) {
      Vue.set(state.names, user.id, user.name)
    }
    for (let id in groups) {
      if (id) {
        Vue.set(state.names, id, groups[id].name)
      }
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
  ownerLoading (state) {
    return state.loading > 0
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
