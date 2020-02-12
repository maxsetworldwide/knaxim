import { shallowMount } from '@vue/test-utils'
import HeaderSearchHistory from '@/components/header-search-history'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(HeaderSearchHistory, {
    propsData: {
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      searchHistory () {
        return []
      },
      ...options.computed
    }
  })
}

describe('HeaderSearchHistory', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(HeaderSearchHistory)).toBe(true)
  })
})
