import { shallowMount } from '@vue/test-utils'
import RegistrationModal from '@/components/modals/registration-modal'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(RegistrationModal, {
    stubs: ['b-img', 'b-form-group', 'b-form-invalid-feedback', 'b-form',
      'b-modal', 'b-form-input', 'b-form-text', 'b-button', 'b-form-checkbox'],
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

describe('RegistrationModal', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(RegistrationModal)).toBe(true)
  })
})
