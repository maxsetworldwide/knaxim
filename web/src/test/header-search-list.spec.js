import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import HeaderSearchList from '@/components/header-search-list'
import HeaderSearchRow from '@/components/header-search-row'
import { SEARCH, SEARCH_TAG } from '@/store/actions.type'
import SearchService from '@/service/search'
import { flushPromises, TestStore } from './utils'

const localVue = createLocalVue()
localVue.use(Vuex)

const testUser = {
  id: 'testID',
  name: 'testName'
}
const testGroup = {
  id: 'testGID',
  name: 'groupName'
}
const testFind = 'testFindPhrase'
const testStore = new TestStore({
  actions: {
    [SEARCH] () {
      return []
    },
    [SEARCH_TAG] () {
      return []
    }
  },
  state: {
    testActiveGroup: null,
    testUser
  },
  getters: {
    searchMatches: () => [],
    searchLines: () => {},
    activeGroup: (state) => state.testActiveGroup,
    currentUser: (state) => state.testUser,
    loading: () => false,
    populateFiles: () => () => {}
  },
  mutations: {
    testSetActiveGroup (state, payload) {
      state.testActiveGroup = payload
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
  return shallowMount(HeaderSearchList, {
    stubs: ['b-container'],
    store,
    localVue,
    propsData: {
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

describe('HeaderSearchList', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(HeaderSearchList)).toBe(true)
  })
  it('dispatches SEARCH', () => {
    let store = testStore.createStore()
    spyOn(store, 'dispatch')
    shallowMountFa({ store })
    expect(store.dispatch).toHaveBeenCalledWith(SEARCH, {
      find: testFind,
      acr: undefined
    })
    const dispatchCalls = store.dispatch.calls.allArgs()
    dispatchCalls.forEach((args) => {
      expect(args[0]).not.toEqual(SEARCH_TAG)
    })
  })
  it('dispatches SEARCH_TAG', () => {
    let store = testStore.createStore()
    spyOn(store, 'dispatch')
    const tag = 'testTag'
    shallowMountFa({ props: { tag }, store })
    const context = SearchService.newOwnerContext(testUser.id)
    const match = SearchService.newMatchCondition(
      testFind,
      tag,
      false,
      testUser.id
    )
    expect(store.dispatch).toHaveBeenCalledWith(SEARCH_TAG, { context, match })
    const dispatchCalls = store.dispatch.calls.allArgs()
    dispatchCalls.forEach((args) => {
      expect(args[0]).not.toEqual(SEARCH)
    })
  })
  it('searches when find changes', async () => {
    let store = testStore.createStore()
    spyOn(store, 'dispatch')
    const wrapper = shallowMountFa({ store })
    const reducer = (acc, args) => (args[0] === SEARCH ? acc + 1 : acc)
    const preDispatchAmount = store.dispatch.calls.allArgs().reduce(reducer, 0)
    const newFind = 'newFindPhrase'
    wrapper.vm.$props.find = newFind
    await flushPromises()
    const postDispatchAmount = store.dispatch.calls.allArgs().reduce(reducer, 0)
    expect(postDispatchAmount).toEqual(preDispatchAmount + 1)
  })
  it('searches when active group changes', async () => {
    let store = testStore.createStore()
    spyOn(store, 'dispatch')
    const wrapper = shallowMountFa({ store })
    const reducer = (acc, args) => (args[0] === SEARCH ? acc + 1 : acc)
    const preDispatchAmount = store.dispatch.calls.allArgs().reduce(reducer, 0)
    wrapper.vm.$store.commit('testSetActiveGroup', testGroup)
    await flushPromises()
    const postDispatchAmount = store.dispatch.calls.allArgs().reduce(reducer, 0)
    expect(postDispatchAmount).toEqual(preDispatchAmount + 1)
  })
  it('passes props to header-search-row', () => {
    const props = {
      find: 'testFind',
      acr: 'TACR'
    }
    const fExt = 'test'
    const fName = 'testFName.' + fExt
    const fid = 'testfid'
    const store = testStore.createStore({
      getters: {
        searchMatches: () => [{ name: fName, id: fid }]
      }
    })
    const wrapper = shallowMountFa({ props, store })
    const rows = wrapper.findAll(HeaderSearchRow)
    expect(rows.length).toEqual(1)
    const expectedRowProps = {
      webpage: false,
      name: fName,
      ext: fExt,
      id: fid,
      find: props.find,
      acr: props.acr
    }
    expect(rows.at(0).vm.$props).toEqual(expectedRowProps)
  })
})
