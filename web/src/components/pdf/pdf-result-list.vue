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
  <b-list-group class="list h-100 w-100">
    <h5 class="text-center">Matches:</h5>
    <b-list-group-item
      flush
      button
      class="py-1 item result-text"
      @click.stop.prevent="handleClick(match)"
      v-for="(match, index) in matchList"
      :key="index"
    >
      <span>Pg.{{ match.page }}:{{ match.sentenceText }}</span>
    </b-list-group-item>
  </b-list-group>
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
}

.result-text {
  font-size: 0.8rem;
  width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.list {
  overflow: auto;
}
</style>
