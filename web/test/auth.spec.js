import { shallowMount } from '@vue/test-utils'
import Auth from '@/components/auth'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html
const $route = {
  name: 'login'
}

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(Auth, {
    stubs: [],
    mocks: {
      $route
    },
    propsData: {
      ...options.props
    },
    methods: {
      openLogin () {
        return true
      },
      openReg () {
        return true
      },
      ...options.methods
    },
    computed: {
      ...options.computed
    }
  })
}

describe('Auth', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(Auth)).toBe(true)
  })
})
