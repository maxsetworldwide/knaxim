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
