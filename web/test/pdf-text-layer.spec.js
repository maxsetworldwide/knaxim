import { shallowMount, createLocalVue } from '@vue/test-utils'
import PdfTextLayer from '@/components/pdf/pdf-text-layer'
import Vuex from 'vuex'

const localVue = createLocalVue()
localVue.use(Vuex)

const highlightString = 'HIGHLIGHTSTRING'
const highlightClassName = 'keyword'
let store = new Vuex.Store({
  getters: {
    currentSearch () {
      return highlightString
    }
  }
})

const testWord = 'weWantThisStringToAppearInTheDocument'
const testTextContent = {
  items: [
    {
      str: `${testWord} and the next will be highlighted: ${highlightString}`,
      dir: 'ltr',
      width: 90,
      height: 10,
      transform: [10, 0, 0, 10, 56.8, 776.789],
      fontName: 'g_d1_f2'
    }
  ],
  styles: {
    g_d1_f2: {
      fontFamily: 'monospace',
      ascent: 0.83251953125,
      descent: -0.30029296875,
      vertical: false
    }
  }
}

const testViewport = {
  viewBox: [0, 0, 595.303937007874, 841.889763779528],
  scale: 3.4963235294117645,
  rotation: 0,
  offsetX: 0,
  offsetY: 0,
  transform: [
    3.4963235294117645,
    0,
    0,
    -3.4963235294117645,
    0,
    2943.5189902732764
  ],
  width: 2081.3751621120887,
  height: 2943.5189902732764
}

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(PdfTextLayer, {
    store,
    localVue,
    propsData: {
      page: {
        getViewport () {
          return testViewport
        },
        getContext () {
          return true
        },
        getTextContent () {
          return new Promise((resolve, reject) => {
            resolve(testTextContent)
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

function flushQueue () {
  return new Promise(resolve => setTimeout(resolve, 0))
}

describe('PdfTextLayer', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(PdfTextLayer)).toBe(true)
  })
  it('emits rendered', async () => {
    const wrapper = shallowMountFa()
    await wrapper.vm.$nextTick()
    expect(wrapper.emitted().rendered).toBeTruthy()
  })
  it('renders provided text', async () => {
    const wrapper = shallowMountFa()
    /* was using wrapper.vm.$nextTick, but required two to work, so using
     * flushQueue() to make this work despite number of ticks required
     */
    await flushQueue()
    expect(wrapper.html()).toContain(testWord)
  })
  it('highlights the correct word', async () => {
    const wrapper = shallowMountFa()
    await flushQueue()
    const expectedString = `<span class="${highlightClassName}">${highlightString}</span>`
    expect(wrapper.html()).toContain(expectedString)
  })
})
