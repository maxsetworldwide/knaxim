<!--
 - nlp-graph: wrapper component for donut-complete for the use of displaying
 -            nlp data from vuex. The parent can choose which tag to display.
 -
 - props:
 -   fid: file id
 -   type: choose which graph to show, specified within the data() field. Any
 -         string within the props field of the type may be used to specify it,
 -         e.g. 't', 'topic', and 'topics' might all be valid to choose Topic.
 - events:
 -   no-data: emitted when the data source of the chosen graph is empty.
 - slots:
 -   default: element appearing above the graph, usually for a title. Falls back
 -            to a capitalized title of the chosen tag in an <h3> tag.
 -->

<template>
  <div>
    <b-spinner v-if="nlpLoading" />
    <div v-else>
      <slot>
        <h3>{{ propToType.title }}</h3>
      </slot>
      <donut-complete
        :dataVals="graphData"
        :colors="propToType.colors"
        @click="handleGraphClick($event)"
      />
    </div>
  </div>
</template>

<script>
import DonutComplete from '@/components/charts/donut-complete'
import { Color } from '@/components/charts/presets'
import { mapGetters, mapActions } from 'vuex'
import { NLP_DATA } from '@/store/actions.type'

export default {
  name: 'nlp-graph',
  components: {
    DonutComplete
  },
  props: {
    fid: {
      type: String,
      required: true
    },
    type: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      types: Object.freeze({
        TOPIC: {
          title: 'Topics',
          tag: 'topic',
          props: ['t', 'topic', 'topics'],
          src: {
            dataLocation: 'nlpTopics',
            index: true // index into dataLocation via this.fid
          },
          colors: Color.Topics
        },
        ACTION: {
          title: 'Actions',
          tag: 'action',
          props: ['a', 'action', 'actions'],
          src: {
            dataLocation: 'nlpActions',
            index: true
          },
          colors: Color.Actions
        },
        RESOURCE: {
          title: 'Resources',
          tag: 'resource',
          props: ['r', 'resource', 'resources'],
          src: {
            dataLocation: 'modifiedResources', // resources rely on topics
            index: false
          },
          colors: Color.Resources
        }
      })
    }
  },
  computed: {
    // translate the given prop to a consistent type to use programmatically
    propToType () {
      const prop = this.type
      for (let key in this.types) {
        let type = this.types[key]
        if (type.props.includes(prop)) {
          return type
        }
      }
      return {
        title: 'Unknown'
      }
    },
    graphSource () {
      const src = this.propToType.src
      if (!src) {
        return []
      } else if (src.index) {
        return this[src.dataLocation][this.fid]
      } else {
        return this[src.dataLocation]
      }
    },
    graphData () {
      let data = this.graphSource || []
      data = data.filter((datum) => {
        return datum.word.length > 2
      }).slice(0, 7)
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
    ...mapGetters(['nlpLoading', 'nlpTopics', 'nlpActions', 'nlpResources']),
    // resources will want to filter out words that exist in topics
    modifiedResources () {
      const topicType = this.types.TOPIC
      const topicSrc = (topicType.src.index ? this[topicType.src.dataLocation][this.fid] : this[topicType.src.dataLocation]) || []
      if (!this.nlpResources[this.fid]) {
        return []
      }
      const topics = topicSrc.map((topic) => {
        return topic.word || ''
      }).filter((word) => {
        return word !== ''
      })
      return this.nlpResources[this.fid].filter(({ word }) => {
        return !topics.includes(word)
      })
    }
  },
  methods: {
    handleGraphClick (label) {
      const tag = this.propToType.tag
      this.$router.push({ path: `/search/${label}/tag/${tag}` })
    },
    ...mapActions([NLP_DATA])
  },
  created () {
    const { fid } = this
    const category = this.propToType.tag
    if (category) {
      this[NLP_DATA]({ fid, category, start: 0, end: 20 }).finally(() => {
        if (!(this.graphSource && this.graphSource.length)) {
          this.$emit('no-data')
        }
      })
    } else {
      this.$emit('no-data')
    }
  }
}
</script>

<style>

</style>
