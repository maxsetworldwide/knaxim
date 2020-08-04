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
import { shallowMount, createLocalVue } from '@vue/test-utils'
import { ACRONYMS } from '@/store/actions.type'
import AcronymSearch from '@/components/acronym-search'
import Vuex from 'vuex'
import { flushPromises, TestStore } from './utils'

const localVue = createLocalVue()
localVue.use(Vuex)

const testPhrase = 'MANPADS'
const testResult = 'Man-Portable Air-Defense System'
const testSlot = '<span>{{ props.result }}</span>'
const expectedHTML = `<span>${testResult}</span>`
const testStore = new TestStore({
  actions: {
    [ACRONYMS] (ctx) {
      ctx.commit('testSetAcronym', testResult)
    }
  },
  state: {
    testResult: ''
  },
  getters: {
    acronymResults: (state) => {
      return state.testResult
    }
  },
  mutations: {
    testSetAcronym (state, payload) {
      state.testResult = payload
    }
  }
})

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (
  options = { props: {}, methods: {}, computed: {}, store: null }
) => {
  let store = options.store || testStore.createStore()
  return shallowMount(AcronymSearch, {
    store,
    localVue,
    scopedSlots: {
      default: testSlot
    },
    propsData: {
      ...options.props,
      phrase: testPhrase
    },
    methods: {
      ...options.methods
    },
    computed: {
      ...options.computed,
      result () {
        return testResult
      }
    }
  })
}

describe('AcronymSearch', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(AcronymSearch)).toBe(true)
  })
  it('slots acronym result', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.html()).toContain(expectedHTML)
  })
  it('dispatches acronym search while typing', async () => {
    const store = testStore.createStore()
    spyOn(store, 'dispatch')
    const wrapper = shallowMountFa({ store })
    const reducer = (acc, args) => (args[0] === ACRONYMS ? acc + 1 : acc)
    const preDispatchAmount = store.dispatch.calls.allArgs().reduce(reducer, 0)
    wrapper.vm.$props.phrase = 'newPhrase'
    await flushPromises()
    const postDispatchAmount = store.dispatch.calls.allArgs().reduce(reducer, 0)
    expect(postDispatchAmount).toBeGreaterThan(preDispatchAmount)
  })
})
