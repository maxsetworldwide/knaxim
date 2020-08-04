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
import Vuex from 'vuex'
import HeaderSearchRow from '@/components/header-search-row'
import FileIcon from '@/components/file-icon'
import { TestStore } from './utils'

const localVue = createLocalVue()
localVue.use(Vuex)

const testName = 'testName'
const testExt = 'tst'
const testFind = 'testFind'
const testID = 'testID'
const testLines = {
  [testID]: {
    loading: false,
    loadingCount: 0,
    matched: [
      {
        Content: ['content line number one'],
        Position: 0
      },
      {
        Content: ['content test this time the second'],
        Position: 1
      }
    ]
  }
}
const testStore = new TestStore({
  actions: {},
  state: {},
  getters: {
    searchLines: () => testLines
  },
  mutations: {}
})

const shallowMountFa = (
  options = { props: {}, methods: {}, computed: {}, store: null }
) => {
  let store = options.store || testStore.createStore()
  return shallowMount(HeaderSearchRow, {
    // stubs: ['b-row', 'b-col'],
    store,
    localVue,
    propsData: {
      name: testName,
      ext: testExt,
      id: testID,
      find: testFind,
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

describe('HeaderSearchRow', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(HeaderSearchRow)).toBe(true)
  })
  it('renders title', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.text()).toContain(testName)
  })
  it('renders matched lines', () => {
    const wrapper = shallowMountFa()
    testLines[testID].matched.forEach((line) => {
      expect(wrapper.text()).toContain(line.Content)
    })
  })
  it('escapes html', () => {
    const contentHTML = '<img src="" onerror="xss code"></img>'
    const lineOverwrite = {
      [testID]: {
        loading: false,
        loadingCount: 0,
        matched: [
          {
            Content: [contentHTML],
            Position: 0
          }
        ]
      }
    }
    const store = testStore.createStore({
      getters: {
        searchLines: () => lineOverwrite
      }
    })
    const wrapper = shallowMountFa({ store })
    expect(wrapper.html()).not.toContain(contentHTML)
    expect(wrapper.text()).toContain(contentHTML)
  })
  it('uses file icon correctly', () => {
    const wrapper = shallowMountFa()
    const icon = wrapper.find(FileIcon)
    const expectedProps = {
      extention: testExt,
      webpage: false
    }
    expect(icon.props()).toEqual(expectedProps)
  })
})
