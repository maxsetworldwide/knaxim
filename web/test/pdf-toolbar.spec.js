import { shallowMount } from '@vue/test-utils'
import PdfToolbar from '@/components/pdf/pdf-toolbar'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(PdfToolbar, {
    stubs: ['b-col', 'b-button', 'b-row'],
    propsData: {
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      isFavorite () {
        return true
      },
      ...options.computed
    }
  })
}

describe('PdfToolbar', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfToolbar)).toBe(true)
  })
})
