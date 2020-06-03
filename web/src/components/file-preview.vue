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
