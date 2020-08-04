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
import { shallowMount } from '@vue/test-utils'
import DonutComplete from '@/components/charts/donut-complete'
import Donut from '@/components/charts/donut'
import ChartLegend from '@/components/charts/chart-legend'

const testLabels = ['First Label', 'Second Label', 'Third Label']
const testData = [10, 5, 20]
const testDataVals = testLabels.map((label, idx) => {
  return { label, data: testData[idx] }
})
const testChart = {
  data: {
    datasets: [
      {
        backgroundColor: ['rgb(10, 0, 0)', 'rgb(0, 10, 0)', 'rgb(0, 0, 10)'],
        data: testData
      }
    ],
    labels: testLabels
  }
}

const shallowMountInst = (
  options = { props: {}, methods: {}, computed: {} }
) => {
  return shallowMount(DonutComplete, {
    propsData: {
      dataVals: testDataVals
    }
  })
}

describe('DonutComplete', () => {
  it('pass dataVals to Donut', () => {
    const wrapper = shallowMountInst()
    const donutWrapper = wrapper.find(Donut)
    expect(donutWrapper.props().dataVals).toEqual(testDataVals)
  })
  it('propagates click events', () => {
    const wrapper = shallowMountInst()
    const eventsToEmit = ['event one', 'event two']
    const donutWrapper = wrapper.find(Donut)
    const legendWrapper = wrapper.find(ChartLegend)
    donutWrapper.vm.$emit('click', eventsToEmit[0])
    legendWrapper.vm.$emit('click', eventsToEmit[1])
    const emittedEvents = wrapper.emitted().click
    expect(emittedEvents.length).toBe(eventsToEmit.length)
    emittedEvents.forEach((e, idx) => {
      let expected = eventsToEmit[idx]
      expect(e.length).toBe(1)
      expect(e[0]).toBe(expected)
    })
  })
  it('passes chart info from donut to legend', () => {
    const wrapper = shallowMountInst()
    const donutWrapper = wrapper.find(Donut)
    const legendWrapper = wrapper.find(ChartLegend)
    donutWrapper.vm.$emit('rendered', testChart)
    wrapper.vm.$nextTick(() => {
      const legendProps = legendWrapper.props()
      expect(legendProps.chart).toEqual(testChart)
    })
  })
})
