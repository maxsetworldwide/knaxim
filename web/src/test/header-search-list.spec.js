import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import HeaderSearchList from '@/components/header-search-list'
import { FILES_SEARCH } from '@/store/actions.type'

const localVue = createLocalVue()
localVue.use(Vuex)

let actions = {
  [FILES_SEARCH] () {
    return []
  }
}
let store = new Vuex.Store({
  actions,
  getters: {
    searchMatches: () => [],
    searchLines: () => {},
    loading: () => false,
    populateFiles: () => () => {}
  }
})

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(HeaderSearchList, {
    stubs: ['b-container'],
    store,
    localVue,
    propsData: {
      find: 'and',
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      activeGroup () {
        return true
      },
      rows () {
        return []
      },
      searchMatches () {
        return []
      },
      ...options.computed
    }
  })
}

describe('HeaderSearchList', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(HeaderSearchList)).toBe(true)
  })
})
