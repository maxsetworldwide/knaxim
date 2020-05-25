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
 * Simple vuex mocking utility. Per spec file, you can create a new TestStore
 * with a default vuex template (eg { getters: ..., actions: ... }), then each
 * individual test can use TestStore.createStore to make a new Vuex Store for
 * that test, including any specific overwrites the test may want to make, which
 * is also in the form of the template.
 */
export class TestStore {
  constructor (template = {}) {
    this.default = template
  }

  createStore (overwrites = {}) {
    return new Vuex.Store(merge(this.default, overwrites))
  }
}
