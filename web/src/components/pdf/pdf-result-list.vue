<!--
pdf-result-list: list of matches that the user can select, notifying the parent
                 upon selection

props:
  matchList: array of matches containing the sentence text. The elements of
             this list are passed to the parent upon selection.

events:
  'select': emitted upon selection of a match, passing the element of the
            given matchList that the selection corresponded to
  'highlight', boolean: emitted upon change of the 'sentenceHighlight' boolean,
                        indicating a desire to toggle sentence highlighting.
-->
<template>
  <div class="list h-100 w-100 d-none d-lg-inline">
    <b-list-group v-if="matchList.length > 0">
      <span class="text-center">
        <!--
           - <b-form-checkbox
           -   v-b-tooltip.hover="{
           -     title: 'Toggle Sentence Highlighting',
           -     placement: 'top',
           -     boundary: 'window'
           -   }"
           -   v-model="sentenceHighlight"
           -   @change="$emit('highlight', $event)"
           -   switch
           -   inline
           -   class="pl-4"
           -   size="lg"
           - />
           -->
        <span class="title">Matches:</span>
      </span>
      <b-list-group-item
        flush
        button
        class="py-1 w-100 item"
        @click.stop.prevent="handleSelect(match)"
        v-for="(match, index) in matchList"
        :key="index"
      >
        <span class="result-text">
          <span>{{ preMatchContext(match) }}</span>
          <span class="phrase">{{ matchPhrase(match) }}</span>
          <span>{{ postMatchContext(match) }}</span>
        </span>
      </b-list-group-item>
    </b-list-group>
    <div v-else class="no-matches h-100 text-center">
      <div class="title my-2">No Matches for:</div>
      <div>{{ currentSearch }}</div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'pdf-result-list',
  props: {
    matchList: {
      type: Array,
      default () {
        return []
      }
    }
  },
  data () {
    return {
      sentenceHighlight: true
    }
  },
  computed: {
    ...mapGetters(['currentSearch'])
  },
  methods: {
    handleSelect (match) {
      this.$emit('select', match)
    },
    preMatchContext (match) {
      const matchStart = match.match.start.global
      const sentenceStart = match.sentence.start.global
      const substringStart = Math.max(matchStart - sentenceStart - 6, 0)
      const substringEnd = matchStart - sentenceStart
      const context = match.sentence.text.substring(
        substringStart,
        substringEnd
      )
      const result = `Pg.${parseInt(match.page) + 1}:${context}`
      return result
    },
    matchPhrase (match) {
      const matchStart = match.match.start.global
      const sentenceStart = match.sentence.start.global
      const matchEnd = match.match.end.global
      const substringStart = matchStart - sentenceStart
      const substringEnd = matchEnd - sentenceStart
      const result = match.sentence.text.substring(substringStart, substringEnd)
      return result
    },
    postMatchContext (match) {
      const matchEnd = match.match.end.global
      const sentenceStart = match.sentence.start.global
      const substringStart = matchEnd - sentenceStart
      const result = match.sentence.text.substring(substringStart)
      return result
    }
  }
}
</script>

<style scoped lang="scss">
.item {
  background-color: $app-bg1;
  -o-transition: 0.5s;
  -ms-transition: 0.5s;
  -moz-transition: 0.5s;
  -webkit-transition: 0.5s;
  transition: 0.5s;
  &:hover {
    background-color: $app-clr2;
  }
  width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.title {
  text-decoration: underline;
  font-weight: 600;
  font-size: 1.2rem;
}

.result-text {
  font-size: 0.9rem;
}

.phrase {
  background-color: $app-clr2;
}

.list {
  overflow: auto;
}

.no-matches {
  border: 2px solid $app-clr2;
  border-radius: 10px;
  background-color: $app-bg1;
}
</style>
