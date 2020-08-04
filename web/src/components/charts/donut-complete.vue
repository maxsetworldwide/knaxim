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
          :colors="colors"
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
    },
    borderColor: String,
    colors: Array
  },
  data () {
    return {
      chart: {}
    }
  }
}
</script>

<style scoped lang="scss"></style>
