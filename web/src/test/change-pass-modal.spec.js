import { shallowMount, createLocalVue } from '@vue/test-utils'
import ChangePassModal from '@/components/modals/change-pass-modal'

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
    isAuthenticated () {
      return true
    }
  }
})

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(ChangePassModal, {
    stubs: ['b-img', 'b-form-group', 'b-form-invalid-feedback', 'b-form',
      'b-modal', 'b-form-input', 'b-form-text', 'b-button'],
    localVue,
    store,
    propsData: {
      id: 'id-abc-123',
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

describe('ChangePassModal', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(ChangePassModal)).toBe(true)
  })
})
