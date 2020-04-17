import { shallowMount } from '@vue/test-utils'
import Donut from '@/components/charts/donut'

const testLabels = ['Donut Label One', 'Donut Label Two', 'Donut Label Three']
const testData = [25, 25, 50]
const testDataVals = testLabels.map((label, idx) => {
  return { label, data: testData[idx] }
})
const testChart = {
  data: {
    datasets: [
      {
        backgroundColor: ['rgb(50, 0, 0)', 'rgb(0, 50, 0)', 'rgb(0, 0, 50)'],
        data: testData
      }
    ],
    labels: testLabels
  }
}

let spyObject = {
  stubbedRenderChart: function () {
    this.$data._chart = testChart // vue-chartjs exposes the chart object on $data
  }
}
const shallowMountInst = (
  options = {
    props: {},
    methods: {
      renderChart: spyObject.stubbedRenderChart
    },
    computed: {}
  }
) => {
  return shallowMount(Donut, {
    ...options,
    propsData: {
      dataVals: testDataVals
    }
  })
}

describe('Donut', () => {
  it('calls renderChart() and passes props to it', () => {
    spyOn(spyObject, 'stubbedRenderChart')
    const spyFunc = spyObject.stubbedRenderChart
    shallowMountInst()
    expect(spyFunc).toHaveBeenCalled()
    const funcCalls = spyFunc.calls
    const args = funcCalls.argsFor(funcCalls.count() - 1)
    const { labels, datasets } = args[0]
    expect(labels).toEqual(testLabels)
    expect(datasets[0].data).toEqual(testData)
  })
  it('emits rendered event with chart data', () => {
    const wrapper = shallowMountInst()
    const events = wrapper.emitted().rendered
    expect(events.length).toBe(1)
    const emittedChart = events[0]
    expect(emittedChart.length).toBe(1)
    expect(emittedChart[0]).toEqual(testChart)
  })
  it('emits click events', () => {
    // click event is hooked up through chart.js options.
    // Grab the function that the component gives to chart.js for the click
    // event, call the function with a mock event, and check that the component
    // propagated the click event up.
    spyOn(spyObject, 'stubbedRenderChart')
    const spyFunc = spyObject.stubbedRenderChart
    const wrapper = shallowMountInst()
    expect(spyFunc).toHaveBeenCalled()
    const funcCalls = spyFunc.calls
    const args = funcCalls.argsFor(funcCalls.count() - 1)
    expect(args.length).toBe(2)
    const { onClick } = args[1]

    const labelIdx = 1
    const expectedLabel = testLabels[1]
    const mockedEvent = [{ _index: labelIdx }]
    // feel free to add a mock point object if needed in the future, currently
    // unused by the component
    onClick(null, mockedEvent)

    const emitted = wrapper.emitted().click
    expect(emitted.length).toBe(1)
    expect(emitted[0].length).toBe(1)
    expect(emitted[0][0]).toBe(expectedLabel)
  })
})
