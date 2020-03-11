import { shallowMount } from '@vue/test-utils'
import AcronymSearch from '@/components/acronym-search'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(AcronymSearch, {
    propsData: {
      ...options.props,
      phrase: 'MANPADS'
    },
    methods: {
      ...options.methods
    },
    computed: {
      ...options.computed,
      result () {
        return 'Man-Portable Air-Defense System'
      }
    }
  })
}

describe('AcronymSearch', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(AcronymSearch)).toBe(true)
  })
})
