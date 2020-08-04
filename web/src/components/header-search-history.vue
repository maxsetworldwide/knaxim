<!--
// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
-->
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
