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
    if (!id) {
      let e = new Error(`LOAD_OWNER: id = ${id}`)
      context.commit(PUSH_ERROR, e)
      throw e
    }
    if (overwrite || !context.state.names[id]) {
      context.commit(OWNER_LOADING, 1)
      context.commit(SET_OWNER_NAME, { id, name: 'loading...' })

      let results = await Promise.allSettled([
        UserService.info({ id }).then(r => r.data.name),
        GroupService.info({ gid: id }).then(r => r.data.name)
      ])

      let name = 'Unknown'
      if (results[0].status === 'fulfilled') {
        name = results[0].value
      } else if (results[1].status === 'fulfilled') {
        name = results[1].value
      } else {
        context.commit(PUSH_ERROR, new Error(`LOAD_OWNER: ${results[0].reason} + ${results[1].reason}`)) // TODO: resolved by creating new handler for looking up owner based on id alone
      }
      context.commit(SET_OWNER_NAME, { id, name })
      context.commit(OWNER_LOADING, -1)
      return name
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
        for (let result of lookedup) {
          if (result.status === 'fulfilled') {
            foundid = result.value
            break
          }
        }
        if (!foundid) {
          throw new Error(`unable to find owner ${name}: ${lookedup[0].reason} ${lookedup[1].reason}`)
        } else {
          commit(SET_OWNER_NAME, {
            id: foundid,
            name
          })
        }
      } catch (e) {
        commit(PUSH_ERROR, e.addDebug('action LOOKUP_OWNER'))
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
