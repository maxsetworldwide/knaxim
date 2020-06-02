import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import FileSideGraphs from '@/components/file-side-graphs'
import NlpGraph from '@/components/charts/nlp-graph'
import { flushPromises } from './utils'

const localVue = createLocalVue()
localVue.use(Vuex)

const testFid = 'testFid'

const getters = {
  nlpLoading () {
    return false
  }
}
let defaultStore = {
  getters
}

const shallowMountFa = (
  options = { props: {}, methods: {}, computed: {}, store: {} }
) => {
  return shallowMount(FileSideGraphs, {
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

describe('FileSideGraphs', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(FileSideGraphs)).toBe(true)
  })
  it('contains the correct graphs', () => {
    const wrapper = shallowMountFa()
    const graphs = wrapper.findAll(NlpGraph)
    // order of the graphs matters.
    const expectedGraphTypes = [['t', 'topic', 'topics'], ['a', 'action', 'actions'], ['r', 'resource', 'resources']]
    expect(graphs.length).toBe(expectedGraphTypes.length)
    expectedGraphTypes.forEach((exp, idx) => {
      let currGraph = graphs.at(idx)
      expect(currGraph.props('fid')).toEqual(testFid)
      expect(exp).toContain(currGraph.props('type'))
    })
  })
  it('renders no graph when no-data is emitted', async () => {
    let wrapper = shallowMountFa()
    wrapper.find(NlpGraph).vm.$emit('no-data')
    await flushPromises()
    expect(wrapper.findAll(NlpGraph).length).toBe(2)
  })
})
