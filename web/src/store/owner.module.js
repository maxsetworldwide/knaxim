import UserService from '@/service/user'
import GroupService from '@/service/group'
import { LOAD_OWNER } from './actions.type'
import { SET_OWNER_NAME } from './mutations.type'

const state = {
  names: {} // map[ownerid]string
}

const actions = {
  async [LOAD_OWNER] (context, { id, overwrite }) {
    if (overwrite || !context.state.names[id]) {
      let response = await UserService.info({ id })
      if (!response) {
        response = await GroupService.info({ id })
      }
      context.commit(SET_OWNER_NAME, { id, name: response.data.name })
    }
  }
}

const mutations = {
  [SET_OWNER_NAME] (state, { id, name }) {
    state.names[id] = name
  }
}

const getters = {
  ownerNames (state) {
    return state.names
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
