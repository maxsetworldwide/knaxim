<!--
pdf-result-list: list of matches that the user can select, notifying the parent
                 upon selection

props:
  matchList: array of matches containing the sentence text. The elements of
             this list are passed to the parent upon selection.

events:
  'select': emitted upon selection of a match, passing the element of the
            given matchList that the selection corresponded to
-->
<template>
  <div class="list h-100 w-100">
    <b-list-group>
      <h5 class="text-center">Matches:</h5>
      <b-list-group-item
        flush
        button
        class="py-1 w-100 item"
        @click.stop.prevent="handleClick(match)"
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
  </div>
</template>

<script>
export default {
  name: 'pdf-result-list',
  props: {
    matchList: Array
  },
  methods: {
    handleClick (match) {
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

.result-text {
  font-size: 0.9rem;
}

.phrase {
  background-color: $app-clr2;
}

.list {
  overflow: auto;
}
</style>
