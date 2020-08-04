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
  <div class="w-100">
    <b-container>
      <b-row v-if="graphsExist.topic">
         <nlp-graph class="w-100" type="topic" :fid="fid" @no-data="graphsExist.topic = false" />
      </b-row>
      <b-row v-if="graphsExist.action">
         <nlp-graph class="w-100" type="action" :fid="fid" @no-data="graphsExist.action = false" />
      </b-row>
      <b-row v-if="graphsExist.resource">
         <nlp-graph class="w-100" type="resource" :fid="fid" @no-data="graphsExist.resource = false" />
      </b-row>
    </b-container>
  </div>
</template>

<script>
import NlpGraph from '@/components/charts/nlp-graph'

export default {
  name: 'file-side-graphs',
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
  }
}
</script>

<style lang="scss" scoped></style>
