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
 - chart-legend: utilize a chart object from chart.js to create a custom legend
 - with click events
 -
 - props:
 -  chart:
 -    chart object from chart.js. Vue-chartjs exposes this as vm.$data._chart
 - events:
 -  click:
 -    emitted upon clicking a legend item, passing the label that was clicked.
 -->
<template>
  <ul>
    <li
      v-for="(datum, idx) in dataSet"
      :key="idx"
      :id="`${labels[idx]}-${idx}`"
      ref="items"
      @click="handleClick(idx)"
    >
      <div class="color-box" :style="{ backgroundColor: colors[idx] }" />
      <span> {{ labels[idx] }} </span>
    </li>
  </ul>
</template>

<script>
export default {
  name: 'chart-legend',
  props: {
    chart: {
      type: Object,
      required: true
    }
  },
  computed: {
    labels () {
      if (this.chart.data) {
        return this.chart.data.labels || []
      }
      return []
    },
    colors () {
      if (this.chart.data) {
        return this.chart.data.datasets[0].backgroundColor || []
      }
      return []
    },
    dataSet () {
      if (this.chart.data) {
        return this.chart.data.datasets[0].data || []
      }
      return []
    }
  },
  methods: {
    handleClick (idx) {
      this.$emit('click', this.labels[idx])
    }
  }
}
</script>

<style lang="scss" scoped>
.color-box {
  width: 10px;
  height: 10px;
  display: inline-block;
  position: relative;
  margin-right: 5px;
}
ul {
  font-size: 0.9em;
  list-style: none;
  margin: 0;
  padding-left: 0.5em;
  height: 100%;
  width: 100%;
  li {
    width: 100%;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    cursor: pointer;
    transition: 0.2s;
    border-radius: 2px;
    padding-left: 2px;
    &:hover {
      background-color: $app-clr2;
    }
  }
}
</style>
