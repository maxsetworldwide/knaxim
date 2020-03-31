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
      return this.dataVals.map((_, idx) => {
        const pallette = this.defaultColorPallette
        return this.defaultColorPallette[idx % pallette.length]
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
    handleClick (point, event) {
      if (event.length) {
        const idx = event[0]._index
        this.$emit('click', this.labels[idx])
      }
    }
  }
}
</script>

<style lang="scss"></style>
