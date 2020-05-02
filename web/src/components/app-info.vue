<template>
  <b-col class="h-100 d-none d-lg-inline" cols="2">
    <b-row class=" d-none d-lg-inline">
      <b-col cols="4" offset-md="4">
        <b-dropdown
          :disabled="selections.length < 2"
          :text="currSelection"
          size="sm"
        >
          <b-dropdown-item
            href="#"
            v-for="item in selections"
            :key="item"
            @click="currSelection = item"
            >{{ item }}</b-dropdown-item
          >
        </b-dropdown>
      </b-col>
    </b-row>
    <b-col v-if="currSelection === 'Searches'" class=" d-none d-lg-inline">
      <div class="sidebar ml-auto">
        <div class="sidebar-search-list">
          <!-- Search History Items  -->
          <b-link
            v-for="item in expandedSearchMatches"
            :key="item.id"
            class="w-100"
            :class="{ 'active-item': item.isActive }"
            :disabled="!!item.isActive"
            :to="`/file/` + item.id"
          >
            <div class="divider">
              {{ item.name }}
            </div>
          </b-link>
        </div>
      </div>
    </b-col>
    <b-col
      v-else-if="currSelection === 'Graphs'"
      class="h-100 d-none d-lg-inline"
      cols="2"
    >
      <div class="sidebar ml-auto">
        <file-side-graphs :fid="$route.params.id" />
      </div>
    </b-col>
  </b-col>
</template>

<script>
import { mapGetters } from 'vuex'
import FileSideGraphs from '@/components/file-side-graphs'

// TODO: move different "sidebars" into different components, eg search list and file graphs, use this component to wrap and switch between them
export default {
  name: 'app-info',
  components: {
    FileSideGraphs
  },
  props: {
    type: {
      type: [String, null],
      default: null
    }
  },
  data () {
    return {
      currSelection: ''
    }
  },
  computed: {
    ...mapGetters(['searchMatches']),
    expandedSearchMatches () {
      return this.searchMatches.map((item) => {
        return {
          ...item,
          isActive: item.id === this.$route.params.id
        }
      })
    },
    searchIsValid () {
      const matches = this.expandedSearchMatches
      return matches && matches.length > 0
    },
    selections () {
      const result = ['Graphs']
      if (this.searchIsValid) {
        result.push('Searches')
      }
      return result
    }
  },
  mounted () {
    this.currSelection = this.searchIsValid ? 'Searches' : 'Graphs'
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
.sidebar-search-list {
  .divider {
    border-bottom: 2px solid gray;
    padding-top: 5px;
  }
  .active-item {
    color: black;
  }
}

.sidebar {
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
}

.sidebar-toggle-switch {
  cursor: pointer;
}
</style>
