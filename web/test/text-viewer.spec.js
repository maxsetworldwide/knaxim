import { shallowMount } from '@vue/test-utils'
import TextViewer from '@/components/text-viewer'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(TextViewer, {
    stubs: ['b-button', 'b-col', 'b-row', 'b-container'],
    mocks: {
      $route: {
        params: {
          id: 'id-abc-123'
        }
      }
    },
    propsData: {
      finalPage: 1,
      ...options.props
    },
    methods: {
      fetchSlice () {
        return {}
      },
      ...options.methods
    },
    computed: {
      loading () {
        return false
      },
      ...options.computed
    }
  })
}

describe('TextViewer', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(TextViewer)).toBe(true)
  })
})
