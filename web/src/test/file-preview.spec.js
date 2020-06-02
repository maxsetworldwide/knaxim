import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import FilePreview from '@/components/file-preview'
import NlpGraph from '@/components/charts/nlp-graph'
import { flushPromises } from './utils'

const localVue = createLocalVue()
localVue.use(Vuex)

const testFid = 'testFid'
const previewLines = ['test preview', 'with some lines']
const getters = {
  filePreview () {
    let result = {}
    const obj = {
      lines: previewLines,
      loading: false,
      loadingCount: 0
    }
    result[testFid] = obj
    return result
  }
}
let defaultStore = {
  getters
}

const shallowMountFa = (
  options = { props: {}, methods: {}, computed: {}, store: {} }
) => {
  return shallowMount(FilePreview, {
    store: new Vuex.Store({
      ...defaultStore,
      ...options.store
    }),
    localVue,
    propsData: {
      fid: testFid,
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

describe('FilePreview', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(FilePreview)).toBe(true)
  })
  it('contains the correct graphs', () => {
    const wrapper = shallowMountFa()
    const graphs = wrapper.findAll(NlpGraph)
    // Order of the graphs matters.
    // Coupled with the props of each type defined in nlp-graph.vue. Unsure of a
    // better way to do this.
    const expectedGraphTypes = [['t', 'topic', 'topics'], ['a', 'action', 'actions'], ['r', 'resource', 'resources']]
    expect(graphs.length).toBe(expectedGraphTypes.length)
    expectedGraphTypes.forEach((exp, idx) => {
      let currGraph = graphs.at(idx)
      expect(currGraph.props('fid')).toEqual(testFid)
      expect(exp).toContain(currGraph.props('type'))
    })
  })
  it('renders the preview lines', () => {
    let wrapper = shallowMountFa()
    const html = wrapper.html()
    previewLines.forEach((line) => {
      expect(html).toContain(line)
    })
  })
  it('renders no graph when no-data is emitted', async () => {
    let wrapper = shallowMountFa()
    wrapper.find(NlpGraph).vm.$emit('no-data')
    await flushPromises()
    expect(wrapper.findAll(NlpGraph).length).toBe(2)
  })
})
