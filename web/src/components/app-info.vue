<template>
  <b-col cols="2">
    <div class="app-info">
      <div
          class="header-search-list w-75 ml-auto">
        <!-- Search History Items  -->
        <router-link v-for="item in searchMatches"
           :key="item.id"
           class="w-100"
           :to="`/file/` + item.id">
           <div class="divider">
             {{ item.name }}
          </div>
        </router-link>
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
    ...mapGetters([
      'searchMatches'
    ])
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
}

</style>
