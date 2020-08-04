// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
