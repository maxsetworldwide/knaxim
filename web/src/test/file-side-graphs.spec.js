import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import FileSideGraphs from '@/components/file-side-graphs'
import DonutComplete from '@/components/charts/donut-complete'

const localVue = createLocalVue()
localVue.use(Vuex)

const testFid = 'testFid'
const nlpArrays = {
  topic: [
    { count: 9, word: 'firstT' },
    { count: 5, word: 'secondT' },
    { count: 3, word: 'thirdT' }
  ],
  action: [
    { count: 12, word: 'firstA' },
    { count: 6, word: 'secondA' },
    { count: 2, word: 'thirdA' }
  ],
  resource: [
    { count: 10, word: 'firstR' },
    { count: 3, word: 'secondR' },
    { count: 1, word: 'thirdR' }
  ]
}

const getters = {
  nlpTopics () {
    let result = {}
    result[testFid] = nlpArrays.topic
    return result
  },
  nlpActions () {
    let result = {}
    result[testFid] = nlpArrays.action
    return result
  },
  nlpResources () {
    let result = {}
    result[testFid] = nlpArrays.resource
    return result
  }
}
let defaultStore = {
  getters
}

const $router = {
  push: function () {}
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
    mocks: {
      $router
    },
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
    const donuts = wrapper.findAll(DonutComplete)
    // order of the graphs matters.
    const expectedGraphs = [
      nlpArrays.topic,
      nlpArrays.action,
      nlpArrays.resource
    ]
    expect(donuts.length).toBe(expectedGraphs.length)
    expectedGraphs.forEach((exp, idx) => {
      let currDonut = donuts.at(idx)
      let currExpectedProps = exp.map((data) => {
        return {
          data: data.count,
          label: data.word
        }
      })
      let currDataValsProp = currDonut.vm.$props.dataVals
      expect(currDataValsProp).toEqual(currExpectedProps)
    })
  })
  it('renders no graph when no data exists', () => {
    const tempGetters = {
      ...getters,
      nlpActions () {
        return {}
      }
    }
    const store = { getters: tempGetters }
    let wrapper = shallowMountFa({ store })
    const donuts = wrapper.findAll(DonutComplete)
    const expectedGraphs = [nlpArrays.topic, nlpArrays.resource]
    expect(donuts.length).toBe(expectedGraphs.length)
  })
  it('pushes router link', () => {
    spyOn($router, 'push')
    const spyFunc = $router.push
    let wrapper = shallowMountFa()
    // first donut should be topic graph
    const emittedLabel = 'label'
    const tag = 'topic'
    wrapper.find(DonutComplete).vm.$emit('click', emittedLabel)
    const expectedArg = {
      path: `/search/${emittedLabel}/tag/${tag}`
    }
    expect(spyFunc).toHaveBeenCalledWith(expectedArg)
  })
})
