<template>
  <b-spinner v-if="loading"></b-spinner>
  <b-container v-else fluid class="header-search-list">
    <b-row>
      <b-col>
        <file-actions :checkedFiles="found"/>
      </b-col>
    </b-row>
    <header-search-row
      v-for="(row, indx) in rows"
      :key="indx"
      :webpage="row.webpage"
      :name="row.name"
      :ext="row.ext"
      :id="row.id"
      :find="find"
      :acr="acr"
    />
  </b-container>
</template>

<script>
import headerSearchRow from '@/components/header-search-row'
import fileActions from '@/components/file-actions'
import { SEARCH } from '@/store/actions.type'
import { mapGetters } from 'vuex'

export default {
  components: {
    headerSearchRow,
    fileActions
  },
  props: {
    find: {
      type: String,
      required: true
    },
    acr: String
  },
  name: 'header-search-list',
  // TODO: Merge the watch and before Mount logic or create a method to handle
  //  both use cases.  One is for first render the other is for prop change.
  watch: {
    find (val) {
      this.$store.dispatch(SEARCH, { find: val, acr: this.acr })
    },
    activeGroup () {
      this.$store.dispatch(SEARCH, { find: this.find, acr: this.acr })
    }
  },
  beforeMount () {
    this.$store.dispatch(SEARCH, { find: this.find, acr: this.acr })
  },
  computed: {
    rows () {
      return this.searchMatches
        .map(file => {
          let splits = file.name.split('.')
          let row = {
            id: file.id,
            name: file.name,
            webpage: file.name.includes('/'),
            ext: splits.length > 1 ? splits[splits.length - 1] : ''
          }
          return row
        })
        .sort((a, b) => {
          const aLen = this.searchLines[a.id].matched.length || 0
          const bLen = this.searchLines[b.id].matched.length || 0
          return bLen - aLen
        })
    },
    found () {
      return this.populateFiles(this.searchMatches.map(f => f.id))
    },
    ...mapGetters(['searchMatches', 'searchLines', 'activeGroup', 'loading', 'populateFiles'])
  }
}
</script>

<style scoped lang="scss">
.header-search-list {
  .line-no {
    max-width: 3em;
    text-align: right;
    padding-right: 0px;
    direction: rtl;
    overflow: hidden;
  }
  ol {
    list-style: none;
  }
  li {
    line-height: 1.2em;
    padding-top: 0.6em;
  }
  .expand {
    width: 25px;
    height: 25px;
  }
  svg {
    width: 50px;
    height: 50px;
  }
  .lite {
    background: lightyellow;
  }
  .unrecognized {
    text-align: left;
  }
}
</style>
