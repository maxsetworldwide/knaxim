import {
  GET_ERROR
} from './actions.type'

import {
  PUSH_ERROR,
  POP_ERROR
} from './mutations.type'

const state = {
  errors: []
}

const getters = {
  availableErrors ({ errors }) {
    return errors.length > 0
  }
}

const mutations = {
  [PUSH_ERROR] ({ errors }, err) {
    errors.push(err)
  },
  [POP_ERROR] ({ errors }, match) {
    if (errors.length > 0) {
      if (match === errors[0]) {
        errors.shift()
      }
    }
  }
}

const actions = {
  async [GET_ERROR] ({ commit, state }) {
    let err = null
    if (state.errors.length > 0) {
      err = state.errors[0]
      commit(POP_ERROR, err)
    }
    return err
  }
}

export default {
  state,
  getters,
  mutations,
  actions
}
