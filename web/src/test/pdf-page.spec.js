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
import PdfPage from '@/components/pdf/pdf-page'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(PdfPage, {
    propsData: {
      page: {
        getViewport () {
          return true
        },
        getContext () {
          return true
        },
        getTextContent () {
          return new Promise((resolve, reject) => {
            resolve(true)
          })
        }
      },
      ...options.props
    },
    methods: {
      updateElementBounds () {
        return true
      },
      setDimStyle () {
        return true
      },
      drawPage () {
        return true
      },
      destroyPage () {
        return true
      },
      destroyRenderTask () {
        return true
      },
      ...options.methods
    },
    computed: {
      canvasID () {
        return true
      },
      canvasAttrs () {
        return {
          width: 1024,
          height: 768,
          style: 'no',
          class: 'pdf-page'
        }
      },
      actualSizeViewport () {
        return {
          height: 0,
          width: 0
        }
      },
      ...options.computed
    }
  })
}

describe('PdfPage', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfPage)).toBe(true)
  })
})
