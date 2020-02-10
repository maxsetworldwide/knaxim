import { shallowMount } from '@vue/test-utils'
import NavBasicAdd from '@/components/nav-basic-add'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(NavBasicAdd, {
    stubs: ['b-dropdown-item', 'b-dropdown-divider', 'b-modal', 'b-dropdown'],
    propsData: {
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      ...options.computed
    },
    directives: {
      bModal: {
        inserted: function (el) {
          return true
        }
      }
    }
  })
}

describe('NavBasicAdd', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(NavBasicAdd)).toBe(true)
  })
})
