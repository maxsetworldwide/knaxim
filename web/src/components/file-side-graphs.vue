<template>
  <div class="w-100">
    <b-spinner v-if="nlpLoading" />
    <b-container v-else>
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
import { mapGetters, mapActions } from 'vuex'
import { NLP_DATA } from '@/store/actions.type'
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
  },
  computed: {
    ...mapGetters(['nlpLoading'])
  },
  methods: {
    ...mapActions([NLP_DATA])
  },
  created () {
    const { fid } = this
    this[NLP_DATA]({ fid, category: 't', start: 0, end: 7 })
    this[NLP_DATA]({ fid, category: 'a', start: 0, end: 7 })
    this[NLP_DATA]({ fid, category: 'r', start: 0, end: 14 })
  }
}
</script>

<style lang="scss" scoped></style>
