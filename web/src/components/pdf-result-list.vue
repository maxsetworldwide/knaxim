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
  <b-list-group class="list h-100">
    <h5 class="text-center">Matches:</h5>
    <b-list-group-item
      flush
      button
      class="py-1 item"
      @click.stop.prevent="handleClick(match)"
      v-for="(match, index) in matchList"
      :key="index"
    >
      <span class="result-text">{{ shortenedSentence(match.sentenceText) }}</span>
    </b-list-group-item>
  </b-list-group>
</template>

<script>
const PREVIEW_LENGTH = 18
export default {
  name: 'pdf-result-list',
  props: {
    matchList: Array
  },
  methods: {
    handleClick (match) {
      // console.log(match)
      this.$emit('select', match)
    },
    shortenedSentence (sentence) {
      const type = typeof sentence
      if (type === 'string') {
        return sentence.substring(0, PREVIEW_LENGTH) + '...'
      } else {
        // console.log('sentence not string: ', sentence, 'type: ', type)
        return ''
      }
    }
  }
}
</script>

<style scoped lang="scss">

.item {
  background-color: $app-bg1;
  -o-transition:.5s;
  -ms-transition:.5s;
  -moz-transition:.5s;
  -webkit-transition:.5s;
  transition:.5s;
  &:hover {
    background-color: $app-clr2;
  }
}

.result-text {
  font-size: .8rem;
}

.list {
  overflow: auto;
}

</style>
