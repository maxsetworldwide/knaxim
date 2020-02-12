import { shallowMount } from '@vue/test-utils'
import MemberList from '@/components/member-list'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(MemberList, {
    stubs: ['b-input', 'b-modal'],
    propsData: {
      members: [],
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      activeGroup () {
        return false
      },
      ...options.computed
    },
    directives: {
      bModal: {
        inserted (el) {
          return true
        }
      }
    }
  })
}

describe('MemberList', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(MemberList)).toBe(true)
  })
})
