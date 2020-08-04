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
import { flushPromises } from './utils'
import PdfTextLayer from '@/components/pdf/pdf-text-layer'
import Vuex from 'vuex'

const localVue = createLocalVue()
localVue.use(Vuex)

const highlightString = 'HIGHLIGHTSTRING'
const highlightClassName = 'keyword'
let store = new Vuex.Store({
  getters: {
    currentSearch () {
      return highlightString
    }
  }
})

const testWord = 'weWantThisStringToAppearInTheDocument'
const testTextContent = {
  items: [
    {
      str: `${testWord}. and the next will be highlighted. Here it is: ${highlightString}`,
      dir: 'ltr',
      width: 90,
      height: 10,
      transform: [10, 0, 0, 10, 56.8, 776.789],
      fontName: 'g_d1_f2'
    }
  ],
  styles: {
    g_d1_f2: {
      fontFamily: 'monospace',
      ascent: 0.83251953125,
      descent: -0.30029296875,
      vertical: false
    }
  }
}

const testViewport = {
  viewBox: [0, 0, 595.303937007874, 841.889763779528],
  scale: 3.4963235294117645,
  rotation: 0,
  offsetX: 0,
  offsetY: 0,
  transform: [
    3.4963235294117645,
    0,
    0,
    -3.4963235294117645,
    0,
    2943.5189902732764
  ],
  width: 2081.3751621120887,
  height: 2943.5189902732764
}

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(PdfTextLayer, {
    store,
    localVue,
    propsData: {
      page: {
        getViewport () {
          return testViewport
        },
        getContext () {
          return true
        },
        getTextContent () {
          return new Promise((resolve, reject) => {
            resolve(testTextContent)
          })
        }
      },
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      ...options.computed
    }
  })
}

describe('PdfTextLayer', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfTextLayer)).toBe(true)
  })
  it('emits rendered', async () => {
    const wrapper = shallowMountFa()
    await flushPromises()
    expect(wrapper.emitted().rendered).toBeTruthy()
  })
  it('renders only once on load', async () => {
    const wrapper = shallowMountFa()
    await flushPromises()
    expect(wrapper.emitted().rendered.length).toBe(1)
  })
  it('renders provided text', async () => {
    const wrapper = shallowMountFa()
    await flushPromises()
    expect(wrapper.html()).toContain(testWord)
  })
  it('highlights the correct word', async () => {
    const wrapper = shallowMountFa()
    await flushPromises()
    const expectedString = `<span class="${highlightClassName}">${highlightString}</span>`
    expect(wrapper.html()).toContain(expectedString)
  })
  it('emits matches', async () => {
    const wrapper = shallowMountFa()
    await flushPromises()
    expect(wrapper.emitted().matches).toBeTruthy()
  })
  it('returns correct matches', async () => {
    const wrapper = shallowMountFa()
    await flushPromises()
    const matches = wrapper.emitted().matches[0][0].matches
    const sentence = matches[0].sentence.text
    expect(matches.length).toEqual(1)
    expect(sentence).toContain(highlightString)
  })
  it('re-renders text when requested, and only once', async () => {
    const wrapper = shallowMountFa()
    await flushPromises()
    const numRenders = wrapper.emitted().rendered.length
    wrapper.vm.refresh()
    await flushPromises()
    expect(wrapper.emitted().rendered.length).toEqual(numRenders + 1)
  })
})
