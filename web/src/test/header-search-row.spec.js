import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import HeaderSearchRow from '@/components/header-search-row'
import FileIcon from '@/components/file-icon'
import merge from 'lodash/merge'

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
const createStore = function (overwrites = {}) {
  const defaultStoreObj = {
    actions: {},
    state: {},
    getters: {
      searchLines: () => testLines
    },
    mutations: {}
  }
  return new Vuex.Store(merge(defaultStoreObj, overwrites))
}

const shallowMountFa = (
  options = { props: {}, methods: {}, computed: {}, store: null }
) => {
  let store = options.store || createStore()
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
    const store = createStore({
      getters: {
        searchLines: () => lineOverwrite
      }
    })
    const wrapper = shallowMountFa({ store })
    expect(wrapper.html()).not.toContain(contentHTML)
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
