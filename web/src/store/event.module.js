// import {
//   EMIT
// } from './actions.type'
import {
  ON,
  OFF,
  EMIT
} from './mutations.type'

const state = {
  handlers: {}
}

// const getters = {
//
// }

// const actions = {
//
// }

const mutations = {
  [ON] ({ handlers }, { evnt, handler }) {
    if (!handlers[evnt]) {
      handlers[evnt] = [handler]
    } else {
      handlers[evnt].push(handler)
    }
  },
  [OFF] ({ handlers }, { evnt, handler }) {
    let removed = false
    handlers[evnt] = handlers[evnt].filter((h) => {
      if (removed) {
        return true
      }
      if (h === handler) {
        removed = true
        return false
      }
      return true
    })
  },
  [EMIT] ({ handlers }, { evnt, payload }) {
    if (handlers[evnt]) {
      handlers[evnt].forEach((h) => {
        h(payload)
      })
    }
  }
}

export default {
  state,
  // actions,
  mutations
  // getters
}
