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
import NavBasic from '@/components/nav-basic'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

import Vuex from 'vuex'

const localVue = createLocalVue()
localVue.use(Vuex)

let store = new Vuex.Store({
  getters: {
    isAuthenticated () {
      return true
    },
    currentUser () {
      return {
        name: 'test'
      }
    },
    activeGroup () {
      return null
    }
  }
})

const $route = {
  name: 'login'
}

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(NavBasic, {
    stubs: ['b-nav-item', 'b-nav-item', 'b-nav'],
    localVue,
    store,
    mocks: {
      $route
    },
    propsData: {
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      cloudtype () {
        return true
      },
      activeGroup () {
        return true
      },
      groupMode () {
        return true
      },
      ...options.computed
    }
  })
}

describe('NavBasic', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(NavBasic)).toBe(true)
  })
})
