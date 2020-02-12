import { shallowMount } from '@vue/test-utils'
import PdfPage from '@/components/pdf/pdf-page'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(PdfPage, {
    propsData: {
      page: {
        getViewport () {
          return true
        },
        getContext () {
          return true
        },
        getTextContent () {
          return new Promise((resolve, reject) => {
            resolve(true)
          })
        }
      },
      ...options.props
    },
    methods: {
      updateElementBounds () {
        return true
      },
      setDimStyle () {
        return true
      },
      findMatches () {
        return true
      },
      sendMatches () {
        return true
      },
      renderText () {
        return true
      },
      highlightMatches () {
        return true
      },
      appendTextChild () {
        return true
      },
      drawPage () {
        return true
      },
      destroyPage () {
        return true
      },
      destroyRenderTask () {
        return true
      },
      ...options.methods
    },
    computed: {
      canvasID () {
        return true
      },
      canvasAttrs () {
        return {
          width: 1024,
          height: 768,
          style: 'no',
          class: 'pdf-page'
        }
      },
      ...options.computed
    }
  })
}

describe('PdfPage', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfPage)).toBe(true)
  })
})
