import { shallowMount } from '@vue/test-utils'
import TeamControl from '@/components/team-control'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(TeamControl, {
    scopedSlots: {
      default: '<div></div>'
    },
    propsData: {
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
    }
  })
}

describe('TeamControl', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(TeamControl)).toBe(true)
  })
})
