import { shallowMount } from '@vue/test-utils'
import AppInfo from '@/components/app-info'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(AppInfo, {
    stubs: ['router-link', 'b-col'],
    propsData: {
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      expandedSearchMatches () {
        return [{
          id: 'id-abc-123',
          name: 'fake-file',
          isActive: true
        }]
      },
      searchMatches () {
        return [{
          id: 'id-abc-123',
          name: 'fake-file'
        }]
      },
      ...options.computed
    }
  })
}

describe('AppInfo', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(AppInfo)).toBe(true)
  })
})
