<!--
 -  donut.vue: a donut chart based on arbitrary data
 -
 -  props:
 -    borderColor:
 -      string representing the color to border and separate each slice.
 -      Defaults to 'rgba(0, 0, 0, 0)'
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
      default: 'rgba(0, 0, 0, 0)'
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
    defaultColors () {
      const pallette = this.defaultColorPallette
      return this.dataVals.map((_, idx) => {
        return pallette[idx % pallette.length]
      })
    }
  },
  data () {
    return {
      defaultColorPallette: [
        '#e41a1c',
        '#377eb8',
        '#4daf4a',
        '#984ea3',
        '#ff7f00',
        '#ffff33',
        '#a65628'
      ],
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
              const label = data.labels[tooltipItem.index] || ''
              const dataset = data.datasets[tooltipItem.datasetIndex].data
              const total = dataset.reduce((acc, val) => {
                return acc + val
              }, 0)
              const curr = dataset[tooltipItem.index]
              const percent = ((curr / total) * 100).toFixed(1) + '%'
              return [`${label}:`, `${percent}`]
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
          backgroundColor: this.defaultColors,
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
