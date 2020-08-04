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
import PdfDoc from '@/components/pdf/pdf-doc'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

import Vuex from 'vuex'
import { GET_FILE } from '@/store/actions.type'

const localVue = createLocalVue()
localVue.use(Vuex)

let store = new Vuex.Store({
  getters: {
    populateFiles (state) {
      return function (id) {
        return {}
      }
    },
    currentSearch () {
      return 'aa'
    }
  },
  actions: {
    [GET_FILE] () {
      return {}
    }
  }
})

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(PdfDoc, {
    stubs: ['b-col', 'b-row', 'b-container'],
    localVue,
    store,
    propsData: {
      ...options.props
    },
    methods: {
      updateScrollBounds () {
        return 0
      },
      ...options.methods
    },
    computed: {
      url () {
        return '<html></html>'
      },
      ...options.computed
    }
  })
}

describe('PdfDoc', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfDoc)).toBe(true)
  })
})
