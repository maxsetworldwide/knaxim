import { shallowMount, createLocalVue } from '@vue/test-utils'
import StorageInfo from '@/components/storage-info'
import { ON, OFF } from '@/store/mutations.type'
import Vuex from 'vuex'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const localVue = createLocalVue()
localVue.use(Vuex)

let store = new Vuex.Store({
  mutations: {
    [ON] () {},
    [OFF] () {}
  }
})

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(StorageInfo, {
    localVue,
    store,
    propsData: {
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      currentUser () {
        return true
      },
      ...options.computed
    }
  })
}

describe('StorageInfo', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(StorageInfo)).toBe(true)
  })
})
