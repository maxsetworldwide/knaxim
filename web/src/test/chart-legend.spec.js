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
import ChartLegend from '@/components/charts/chart-legend'

const testChart = {
  data: {
    datasets: [
      {
        backgroundColor: ['rgb(255, 0, 0)', 'rgb(0, 255, 0)', 'rgb(0, 0, 255)'],
        data: [5, 5, 10]
      }
    ],
    labels: ['Label One', 'Label Two', 'Label Three']
  }
}
const shallowMountInst = (
  options = { props: {}, methods: {}, computed: {} }
) => {
  return shallowMount(ChartLegend, {
    propsData: {
      chart: testChart
    }
  })
}

describe('ChartLegend', () => {
  it('lists all labels', () => {
    const wrapper = shallowMountInst()
    const items = wrapper.findAll('li')
    expect(items.length).toBe(testChart.data.labels.length)
    testChart.data.labels.forEach((label, idx) => {
      let currItem = items.at(idx)
      expect(currItem.text()).toBe(label)
    })
  })
  it('displays the correct color for each label', () => {
    const wrapper = shallowMountInst()
    const items = wrapper.findAll('li')
    testChart.data.datasets[0].backgroundColor.forEach((color, idx) => {
      let currListItem = items.at(idx)
      let div = currListItem.find('div')
      expect(div.element.style.backgroundColor).toBe(color)
    })
  })
  it('emits click events for each label', () => {
    const wrapper = shallowMountInst()
    const items = wrapper.findAll('li')
    const labels = testChart.data.labels
    labels.forEach((label, idx) => {
      let currItem = items.at(idx)
      currItem.trigger('click')
    })
    const events = wrapper.emitted().click
    expect(events.length).toBe(labels.length)
    events.forEach((e, idx) => {
      expect(e[0]).toBe(labels[idx])
    })
  })
})
