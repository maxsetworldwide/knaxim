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
<template>
  <b-container>
    <b-row align-h="around">
      <b-col v-if="graphsExist.topic" cols="3">
         <nlp-graph type="topic" :fid="fid" @no-data="graphsExist.topic = false"/>
      </b-col>
      <b-col v-if="graphsExist.action" cols="3">
         <nlp-graph type="action" :fid="fid" @no-data="graphsExist.action = false"/>
      </b-col>
      <b-col v-if="graphsExist.resource" cols="3">
         <nlp-graph type="resource" :fid="fid" @no-data="graphsExist.resource = false"/>
      </b-col>
    </b-row>
    <b-row align-h="around">
      <b-col cols="8" align-self="center">
        <!-- summary -->
        <div class="summary">{{ summary }}</div>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import { mapGetters } from 'vuex'
import NlpGraph from '@/components/charts/nlp-graph'

export default {
  name: 'file-preview',
  components: {
    NlpGraph
  },
  props: {
    fid: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      graphsExist: {
        topic: true,
        action: true,
        resource: true
      }
    }
  },
  computed: {
    summary () {
      return this.filePreview[this.fid].lines
        ? this.filePreview[this.fid].lines.join(' ')
        : ''
    },
    ...mapGetters(['filePreview'])
  }
}
</script>

<style lang="scss" scoped>
.summary {
  --lh: 1.2rem;
  --max-lines: 7;
  line-height: var(--lh);
  max-height: calc(var(--lh) * var(--max-lines));
  overflow: hidden;
}
</style>
