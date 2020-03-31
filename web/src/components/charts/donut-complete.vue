<template>
  <b-container class="m-0 p-0">
    <b-row no-gutters align-v="center">
      <b-col cols="6">
        <chart-legend :chart="chart" @click="handleClick" />
      </b-col>
      <b-col cols="6">
        <donut
          class="chart"
          :dataVals="dataVals"
          @rendered="chart = $event"
          @click="handleClick"
        />
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import Donut from '@/components/charts/donut'
import ChartLegend from '@/components/charts/chart-legend'

export default {
  name: 'donut-complete',
  components: {
    Donut,
    ChartLegend
  },
  props: {
    dataVals: {
      type: Array,
      required: true
    }
  },
  data () {
    return {
      chart: {}
    }
  },
  methods: {
    handleClick (payload) {
      console.log({ payload })
    }
  }
}
</script>

<style scoped lang="scss">
.legend {
}

.chart {
}
</style>
