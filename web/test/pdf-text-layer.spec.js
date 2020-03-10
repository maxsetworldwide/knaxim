import { shallowMount } from '@vue/test-utils'
import PdfTextLayer from '@/components/pdf/pdf-text-layer'

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(PdfTextLayer, {
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
      ...options.methods
    },
    computed: {
      ...options.computed
    }
  })
}

describe('PdfTextLayer', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfTextLayer)).toBe(true)
  })
})
