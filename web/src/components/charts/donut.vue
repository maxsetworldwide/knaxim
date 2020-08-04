<!--
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
-->
<!--
 -  donut.vue: a donut chart based on arbitrary data
 -
 -  props:
 -    borderColor:
 -      string representing the color to border and separate each slice.
 -      Defaults to 'rgba(0, 0, 0, 0)'
 -    colors:
 -      array of strings representing the slice colors to cycle through.
 -      There are default contrasting colors if this property is not defined.
 -    dataVals:
 -      array of objects representing the data to pass to the chart
 -    dataVals[i].label:
 -      string representing label or name of this data point
 -    dataVals[i].data:
 -      number representing the magnitude of this data point. These do not have
 -      to be percentages, as chartjs automatically sums all data and converts
 -      them into percentages.
 -
 -  events:
 -    rendered:
 -      emitted upon rendering the chart, passing the chart data from the render
 -    click:
 -      emitted upon clicking a slice of the chart, passing the label of the
 -      slice that was clicked
 -->
<script>
import { Doughnut } from 'vue-chartjs'

export default {
  name: 'donut',
  mixins: [Doughnut],
  props: {
    borderColor: {
      type: String,
      default: '#FFFFFF'
    },
    colors: {
      type: Array,
      default: () => [
        '#e41a1c',
        '#377eb8',
        '#4daf4a',
        '#984ea3',
        '#ff7f00',
        '#ffff33',
        '#a65628'
      ]
    },
    dataVals: {
      type: Array,
      required: true
    }
  },
  computed: {
    labels () {
      return this.dataVals.map((val) => val.label || '')
    },
    dataPoints () {
      return this.dataVals.map((val) => val.data || 0)
    },
    colorsToSize () {
      const pallette = this.colors
      return this.dataVals.map((_, idx) => {
        return pallette[idx % pallette.length]
      })
    }
  },
  data () {
    return {
      options: {
        responsive: true,
        maintainAspectRatio: true,
        layout: {
          padding: 0
        },
        legend: {
          display: false
        },
        tooltips: {
          enabled: true,
          position: 'nearest',
          displayColors: false,
          callbacks: {
            label: function (tooltipItem, data) {
              return data.labels[tooltipItem.index] || ''
            }
          }
        },
        onClick: this.handleClick,
        hover: {
          // make slice hover change cursor to pointer
          onHover: function (e) {
            const point = this.getElementAtEvent(e)
            if (point.length) {
              e.target.style.cursor = 'pointer'
            } else {
              e.target.style.cursor = 'default'
            }
          }
        }
      }
    }
  },
  mounted () {
    const chartData = {
      labels: this.labels,
      datasets: [
        {
          backgroundColor: this.colorsToSize,
          borderColor: this.borderColor,
          data: this.dataPoints
        }
      ]
    }
    this.renderChart(chartData, this.options)
    this.$emit('rendered', this.$data._chart)
  },
  methods: {
    handleClick (_, event) {
      if (event.length) {
        const idx = event[0]._index
        this.$emit('click', this.labels[idx])
      }
    }
  }
}
</script>

<style lang="scss"></style>
