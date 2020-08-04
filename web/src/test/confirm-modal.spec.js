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
import { shallowMount } from '@vue/test-utils'
import ConfirmModal from '@/components/modals/confirm-modal'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(ConfirmModal, {
    stubs: ['b-modal'],
    propsData: {
      ...options.props
    },
    methods: {
      show () {
        return true
      },
      ...options.methods
    },
    computed: {
      ...options.computed
    }
  })
}

describe('ConfirmModal', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(ConfirmModal)).toBe(true)
  })
})
