<!--
 - donut-complete: a wrapper for donut.vue and chart-legend.vue
 -
 - props:
 -  dataVals:
 -    data to pass directly to the donut. See donut.vue.
 -
 - events:
 -  click:
 -    emitted from both the donut and legend components. See donut.vue and
 -    chart-legend.vue.
 -->
<template>
  <b-container class="m-0 p-0">
    <b-row no-gutters align-v="center">
      <b-col cols="6">
        <chart-legend :chart="chart" @click="$emit('click', $event)" />
      </b-col>
      <b-col cols="6">
        <donut
          class="chart"
          :dataVals="dataVals"
          @rendered="chart = $event"
          @click="$emit('click', $event)"
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
  }
}
</script>

<style scoped lang="scss"></style>
