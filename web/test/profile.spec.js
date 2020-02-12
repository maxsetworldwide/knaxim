import { shallowMount } from '@vue/test-utils'
import Profile from '@/components/profile'

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
  return shallowMount(Profile, {
    stubs: [],
    mocks: {
      $route
    },
    propsData: {
      ...options.props
    },
    methods: {
      showChangePass () {
        return true
      },
      ...options.methods
    },
    computed: {
      ...options.computed
    }
  })
}

describe('Profile', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(Profile)).toBe(true)
  })
})
