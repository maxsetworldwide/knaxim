<template>
  <b-spinner v-if="loading"></b-spinner>
  <b-container v-else fluid class="header-search-list">
    <b-row>
      <b-col>
        <file-actions :checkedFiles="fileObjects" />
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
import SearchService from '@/service/search'
import { SEARCH, SEARCH_TAG } from '@/store/actions.type'
import { SET_ACTIVE_FILES } from '@/store/mutations.type'
import { mapGetters, mapMutations } from 'vuex'

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
    acr: String,
    tag: String
  },
  name: 'header-search-list',
  watch: {
    find (val) {
      this.search()
    },
    activeGroup () {
      this.search()
    },
    searchMatches (n) {
      this.currentFiles(n.map(f => f.id))
    }
  },
  beforeMount () {
    this.search()
  },
  methods: {
    search () {
      if (this.tag) {
        const owner =
          (this.activeGroup && this.activeGroup.id) || this.currentUser.id
        const context = SearchService.newOwnerContext(owner)
        const match = SearchService.newMatchCondition(
          this.find,
          this.tag,
          false,
          owner
        )
        this.$store.dispatch(SEARCH_TAG, { context, match })
      } else {
        this.$store.dispatch(SEARCH, { find: this.find, acr: this.acr })
      }
    },
    ...mapMutations({
      currentFiles: SET_ACTIVE_FILES
    })
  },
  computed: {
    rows () {
      return this.searchMatches
        .map((file) => {
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
    fileObjects () {
      return this.populateFiles(this.searchMatches.map((f) => f.id))
    },
    ...mapGetters([
      'searchMatches',
      'searchLines',
      'activeGroup',
      'currentUser',
      'loading',
      'populateFiles'
    ])
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
