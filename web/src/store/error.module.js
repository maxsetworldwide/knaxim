import {
  GET_ERROR,
  ERROR_LOOP
} from './actions.type'

import {
  PUSH_ERROR,
  POP_ERROR,
  ADD_ERROR_LOOP,
  RESET_ERROR
} from './mutations.type'

const state = {
  errors: [],
  errorLoop: Promise.resolve(true)
}

const getters = {
  availableErrors ({ errors }) {
    return errors.length > 0
  }
}

const mutations = {
  [PUSH_ERROR] ({ errors }, err) {
    // if (!(err instanceof Error)) {
    //   err = new Error(`${err}`)
    // }
    errors.push(err)
  },
  [POP_ERROR] ({ errors }, match) {
    if (errors.length > 0) {
      if (match === errors[0]) {
        errors.shift()
      }
    }
  },
  [ADD_ERROR_LOOP] (state, next) {
    state.errorLoop = state.errorLoop.then(next)
  },
  [RESET_ERROR] (state) {
    state.errors = []
    state.errorLoop = Promise.resolve(true)
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
  },
  [ERROR_LOOP] ({ commit, dispatch, getters }, callback) {
    commit(ADD_ERROR_LOOP, async () => {
      try {
        while (getters.availableErrors) {
          let e = await dispatch(GET_ERROR)
          await callback(e)
        }
        return true
      } catch {
        commit(RESET_ERROR)
      }
      return false
    })
  }
}

export default {
  state,
  getters,
  mutations,
  actions
}
