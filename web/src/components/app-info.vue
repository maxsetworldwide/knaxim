<template>
  <b-col cols="2">
    <div class="app-info d-none d-md-inline">
      <div class="header-search-list w-75 ml-auto">
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
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'app-info',
  props: {
    type: {
      type: [String, null],
      default: null
    }
  },
  computed: {
    ...mapGetters(['searchMatches']),
    expandedSearchMatches () {
      return this.searchMatches.map(item => {
        return {
          ...item,
          isActive: item.id === this.$route.params.id
        }
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
.app-info {
  .divider {
    border-bottom: 2px solid gray;
    padding-top: 5px;
  }
  .active-item {
    color: black;
  }
}
</style>
