<template>
  <b-container>
    <b-row align-h="around">
      <b-col v-if="topicData.length > 0" cols="3">
        <h3>Topics</h3>
        <donut-complete
          :dataVals="topicData"
          @click="handleGraphClick('topic', $event)"
        />
      </b-col>
      <b-col v-if="actionData.length > 0" cols="3">
        <h3>Actions</h3>
        <donut-complete
          :dataVals="actionData"
          @click="handleGraphClick('action', $event)"
        />
      </b-col>
      <b-col v-if="resourceData.length > 0" cols="3">
        <h3>Resources</h3>
        <donut-complete
          :dataVals="resourceData"
          @click="handleGraphClick('resource', $event)"
        />
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
import donutComplete from '@/components/charts/donut-complete'
import { mapGetters } from 'vuex'

export default {
  name: 'file-preview',
  components: {
    donutComplete
  },
  props: {
    fid: {
      type: String,
      required: true
    }
  },
  computed: {
    summary () {
      return this.filePreview[this.fid].lines
        ? this.filePreview[this.fid].lines.join(' ')
        : ''
    },
    // TODO: move this logic to a renderless component
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
    ...mapGetters(['filePreview', 'nlpTopics', 'nlpActions', 'nlpResources'])
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
    }
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
