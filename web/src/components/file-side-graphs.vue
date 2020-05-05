<template>
  <div class="w-100">
    <b-spinner v-if="nlpLoading" />
    <b-container v-else class="donut">
      <b-row v-if="topicData.length > 0">
        <h3>Topics</h3>
        <donut-complete
          class="w-100"
          :dataVals="topicData"
          :colors="topicColors"
          @click="handleGraphClick('topic', $event)"
        />
      </b-row>
      <b-row v-if="actionData.length > 0">
        <h3>Actions</h3>
        <donut-complete
          class="w-100"
          :dataVals="actionData"
          :colors="actionsColors"
          @click="handleGraphClick('action', $event)"
        />
      </b-row>
      <b-row v-if="resourceData.length > 0">
        <h3>Resources</h3>
        <donut-complete
          class="w-100"
          :dataVals="resourceData"
          :colors="resourceColors"
          @click="handleGraphClick('resource', $event)"
        />
      </b-row>
    </b-container>
  </div>
</template>

<script>
import DonutComplete from '@/components/charts/donut-complete'
import { Color } from '@/components/charts/presets'
import { mapGetters, mapActions } from 'vuex'
import { NLP_DATA } from '@/store/actions.type'

export default {
  name: 'file-side-graphs',
  components: {
    DonutComplete
  },
  props: {
    fid: {
      type: String,
      required: true
    }
  },
  computed: {
    // TODO: this is repeated logic from file-preview. Move this to a renderless component.
    topicColors () {
      return Color.Topics
    },
    actionsColors () {
      return Color.Actions
    },
    resourceColors () {
      return Color.Resources
    },
    topicData () {
      return this.buildGraphData(this.nlpTopics[this.fid])
    },
    actionData () {
      return this.buildGraphData(this.nlpActions[this.fid])
    },
    resourceData () {
      const topics = this.nlpTopics[this.fid]
        .map((topic) => {
          return topic.word || ''
        })
        .filter((word) => {
          return word !== ''
        })
      return this.buildGraphData(
        this.nlpResources[this.fid]
          .filter(({ word }) => {
            return !topics.includes(word)
          })
          .slice(0, 7)
      )
    },
    ...mapGetters(['nlpTopics', 'nlpActions', 'nlpResources', 'nlpLoading'])
  },
  methods: {
    buildGraphData (data) {
      let result = []
      for (let val in data) {
        let { word, count } = data[val]
        result.push({
          label: word,
          data: count
        })
      }
      return result
    },
    handleGraphClick (tag, label) {
      this.$router.push({ path: `/search/${label}/tag/${tag}` })
    },
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
