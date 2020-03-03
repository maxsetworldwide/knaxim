import { shallowMount, createLocalVue } from '@vue/test-utils'
import PdfResultList from '@/components/pdf/pdf-result-list'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

import Vuex from 'vuex'
const localVue = createLocalVue()
localVue.use(Vuex)

let store = new Vuex.Store({
  getters: {
    currentSearch () {
      return 'aa'
    }
  }
})
const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(PdfResultList, {
    stubs: ['b-list-group'],
    localVue,
    store,
    propsData: {
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

describe('PdfResultList', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfResultList)).toBe(true)
  })
})
