import { shallowMount } from '@vue/test-utils'
import BatchDelete from '@/components/batch-delete'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(BatchDelete, {
    scopedSlots: {
      default: '<div></div>'
    },
    propsData: {
      files: [{
        id: 'id-abc-123',
        name: 'fakeFile',
        own: true,
        date: {
          upload: '00-00-0000'
        }
      }],
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      ...options.computed
    }
  })
}

describe('BatchDelete', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(BatchDelete)).toBe(true)
  })
})
