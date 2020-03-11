import { shallowMount } from '@vue/test-utils'
import ConfirmModal from '@/components/modals/confirm-modal'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(ConfirmModal, {
    stubs: ['b-modal'],
    propsData: {
      ...options.props
    },
    methods: {
      show () {
        return true
      },
      ...options.methods
    },
    computed: {
      ...options.computed
    }
  })
}

describe('ConfirmModal', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(ConfirmModal)).toBe(true)
  })
})
