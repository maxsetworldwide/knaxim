import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import PdfToolbar from '@/components/pdf/pdf-toolbar'
import { TestStore } from './utils'

const localVue = createLocalVue()
localVue.use(Vuex)

const testStore = new TestStore({
  getters: {
    activeFiles () {
      return []
    }
  }
})

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {}, store: null }) => {
  let store = options.store || testStore.createStore()
  return shallowMount(PdfToolbar, {
    stubs: ['b-col', 'b-button', 'b-row'],
    localVue,
    store,
    propsData: {
      file: {
        name: 'name'
      },
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      isFavorite () {
        return true
      },
      ...options.computed
    }
  })
}

describe('PdfToolbar', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfToolbar)).toBe(true)
  })
})
