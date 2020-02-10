import { shallowMount } from '@vue/test-utils'
import NavBasic from '@/components/nav-basic'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(NavBasic, {
    stubs: ['b-nav-item', 'b-nav-item', 'b-nav'],
    propsData: {
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      cloudtype () {
        return true
      },
      activeGroup () {
        return true
      },
      groupMode () {
        return true
      },
      ...options.computed
    }
  })
}

describe('NavBasic', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(NavBasic)).toBe(true)
  })
})
