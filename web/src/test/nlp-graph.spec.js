// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import NlpGraph from '@/components/charts/nlp-graph'
import DonutComplete from '@/components/charts/donut-complete'
import { NLP_DATA } from '@/store/actions.type'
import { flushPromises, TestStore } from './utils'

const localVue = createLocalVue()
localVue.use(Vuex)

const testFid = 'testfid'

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
const testStore = new TestStore({
  getters: {
    nlpTopics () {
      return {
        [testFid]: nlpArrays.topic
      }
    },
    nlpActions () {
      return {
        [testFid]: nlpArrays.action
      }
    },
    nlpResources () {
      return {
        [testFid]: nlpArrays.resource
      }
    },
    nlpLoading () {
      return false
    }
  },
  actions: {
    async [NLP_DATA] () {
    }
  }
})

const $router = {
  push: function () {}
}

const shallowMountFa = (
  options = { props: {}, methods: {}, computed: {}, store: null }
) => {
  let store = options.store || testStore.createStore()
  return shallowMount(NlpGraph, {
    store,
    localVue,
    mocks: {
      $router
    },
    propsData: {
      fid: testFid,
      type: 'topic',
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

describe('NlpGraph', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(NlpGraph)).toBe(true)
  })
  it('dispatches NLP_DATA for the given type', () => {
    const types = ['action', 'topic', 'resource']
    types.forEach((type) => {
      const store = testStore.createStore()
      spyOn(store, 'dispatch').and.callThrough()
      shallowMountFa({ store, props: { type } })
      const calls = store.dispatch.calls.allArgs().map((call) => {
        let result = { action: call[0] }
        if (call.length > 1) {
          const { fid, category } = call[1]
          result = {
            action: call[0],
            fid,
            category
          }
        }
        return result
      })
      expect(calls).toContain({
        action: NLP_DATA,
        fid: testFid,
        category: type
      })
    })
  })
  it('passes the correct dataVals to donut-complete', () => {
    const wrapper = shallowMountFa({ props: { type: 'action' } })
    const expectedVals = nlpArrays.action.map((data) => {
      return {
        data: data.count,
        label: data.word
      }
    })
    expect(wrapper.find(DonutComplete).props('dataVals')).toEqual(expectedVals)
  })
  it('renders the correct graph based on the type prop and renders a default title', () => {
    const expected = [
      { props: ['a', 'action', 'actions'], title: 'Actions' },
      { props: ['r', 'resource', 'resources'], title: 'Resources' },
      { props: ['t', 'topic', 'topics'], title: 'Topics' }
    ]
    expected.forEach((exp) => {
      exp.props.forEach((type) => {
        let wrapper = shallowMountFa({ props: { type } })
        expect(wrapper.text()).toContain(exp.title)
      })
    })
  })
  it('pushes a router link when graph is clicked', () => {
    spyOn($router, 'push')
    const spyFunc = $router.push
    let wrapper = shallowMountFa()
    const emittedLabel = 'label'
    const tag = 'topic'
    wrapper.find(DonutComplete).vm.$emit('click', emittedLabel)
    const expectedArg = {
      path: `/search/${emittedLabel}/tag/${tag}`
    }
    expect(spyFunc).toHaveBeenCalledWith(expectedArg)
  })
  it('emits no-data when data is missing for chosen graph', async () => {
    const store = testStore.createStore({
      getters: {
        nlpActions () {
          return {}
        }
      }
    })
    const wrapper = shallowMountFa({ store, props: { type: 'action' } })
    await flushPromises()
    expect(wrapper.emitted('no-data').length).toEqual(1)
  })
  it('does not emit no-data when data is missing for a non chosen graph', () => {
    const store = testStore.createStore({
      getters: {
        nlpActions () {
          return {}
        }
      }
    })
    const wrapper = shallowMountFa({ store, props: { type: 'topic' } })
    expect(wrapper.emitted('no-data')).toBeUndefined()
  })
})
