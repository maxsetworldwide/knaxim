import Vue from 'vue'
import UserService from '@/service/user'
import GroupService from '@/service/group'
import { LOAD_OWNER, LOOKUP_OWNER } from './actions.type'
import { SET_OWNER_NAME, PROCESS_SERVER_STATE, OWNER_LOADING, PUSH_ERROR } from './mutations.type'

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
        context.commit(SET_OWNER_NAME, { id, name: response.data.name })
      } catch (err) {
        try {
          response = await GroupService.info({ id })
          context.commit(SET_OWNER_NAME, { id, name: response.data.name })
        } catch (errr) {
          context.commit(PUSH_ERROR, new Error(`LOAD_OWNER: ${err} + ${errr}`))
          context.commit(SET_OWNER_NAME, { id, name: 'Unknown' })
          return 'Unknown'
        }
      } finally {
        context.commit(OWNER_LOADING, -1)
      }
      return response.data.name
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
        let lookedup = await Promise.allSettled([
          UserService.lookup({ name }).then(res => res.data.id),
          GroupService.lookup({ name }).then(res => res.data.id)
        ])
        for (let result in lookedup) {
          if (result.status === 'fulfilled') {
            foundid = result.value
            break
          }
        }
        if (foundid) {
          commit(SET_OWNER_NAME, {
            id: foundid,
            name
          })
        } else {
          throw new Error(`unable to find owner: ${name}`)
        }
      } catch (e) {
        commit(PUSH_ERROR, new Error(`LOOKUP_OWNER: ${e}`))
      } finally {
        commit(OWNER_LOADING, -1)
      }
    }
    return foundid
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
