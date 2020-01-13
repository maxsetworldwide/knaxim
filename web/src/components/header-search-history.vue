<template>
  <div class="header-search-history">
    <!-- Search History Items  -->
    <b-button v-for="item in items"
        :key="item.id"
        class="w-100 m-1"
        :to="`/search/` + item"
        @click="setSearch(item)">
      {{ item }}
    </b-button>
  </div>
</template>

<script>
import { mapMutations, mapGetters } from 'vuex'

export default {
  name: 'header-search-history',
  data () {
    return {
      selected: '',
      search: []
    }
  },

  computed: {
    items () {
      return this.searchHistory.slice(0, 4)
    },
    ...mapGetters([
      'searchHistory'
    ])
  },

  methods: {
    ...mapMutations([
      'doSearch',
      'setSearch'
    ])
  },
  props: {
  },
  components: {
  }
}
</script>

<style lang="scss">
.header-search-history {
}

</style>
