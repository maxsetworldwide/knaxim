import Vuex from 'vuex'
import merge from 'lodash/merge'

/*
 * Useful for awaiting components to change upon data change.
 *
 * Was using wrapper.vm.$nextTick to await DOM changes, but required two+ to
 * work. Awaiting on this function instead seems to fix this.
 */
export function flushPromises () {
  return new Promise((resolve) => setTimeout(resolve, 0))
}

/*
 *
 */
export class TestStore {
  constructor (template = {}) {
    this.default = template
  }

  createStore (overwrites = {}) {
    return new Vuex.Store(merge(this.default, overwrites))
  }
}
