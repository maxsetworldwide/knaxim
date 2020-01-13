<template>
 <b-navbar-nav class="header-search w-100 ml-auto">
  <!-- Search Bar -->
  <!-- TODO: Should the find search bar be updated when a search history is
    clicked?  Then, Fix circular refernce issue bewteen state, and model.
    for the find parameter.
  -->
  <!-- Hide the list after enter is pressed but, only
    after the search is initiated. -->
  <b-form-input ref="search" placeholder="Discover" size="md"
      class="search--input mr-sm-2 w-100"
      v-model="find"
      @focus="onFocus"
      @blur="showList = !showList"
      @keydown.enter="search"
      @keyup.enter="$refs.search.blur"
  ></b-form-input>

  <div class="custom-dropdown w-100" v-if="showList">
    <b-list-group flush>
      <acronym-search ref="acronym" :phrase="find" #default="{ result }">
        <!-- mousedown happens Before the blur event! -->
        <b-list-group-item button class="w-100"
            v-for="(item, index) in result.slice(0, 6)"
            :key="index"
            :value="item"
            @mousedown="searchWithAcronym(item)"
        > {{ item }} </b-list-group-item>

        <!-- Divider -->
        <b-list-group-item v-if="dropDown.length">
          <hr class="w-100">
        </b-list-group-item>
      </acronym-search>

      <!-- TODO: remove getter param in getSearch, just mapState history. -->
      <!-- Search History Items  -->
      <b-list-group-item button class="w-100"
          v-for="item in dropDown"
          :key="item"
          :value="item"
          @mousedown="search(item)"
      > {{ item }} </b-list-group-item>
    </b-list-group>
  </div>
 </b-navbar-nav>
</template>

<script>
import { mapGetters } from 'vuex'
import AcronymSearch from '@/components/acronym-search'

export default {
  name: 'header-search',
  components: {
    AcronymSearch
  },
  data () {
    return {
      find: '',
      showList: false
    }
  },

  computed: {
    dropDown () {
      return this.searchHistory.slice(0, 6)
    },
    ...mapGetters([
      'searchHistory'
    ])
  },

  methods: {
    onFocus () {
      this.showList = !this.showList
      // acronym ref is not available until afer it is visible
      this.$nextTick(() => {
        this.$refs.acronym.search()
      })
    },
    searchWithAcronym (item) {
      let ac = this.find
      if (typeof item === 'string' && item.length > 0) {
        this.find = item
      }
      if (this.find) {
        this.$router.push(({ path: `/search/${this.find}/acronym/${ac}` }))
      }
    },
    search (item) {
      // Load search view
      if (typeof item === 'string' && item.length > 0) {
        this.find = item
      }
      if (this.find) {
        this.$router.push(({ path: `/search/${this.find}` }))
      }
    }
  }
}
</script>

<style lang="scss">
.header-search {
  position: relative;

  // Search container and field.
  .search--input {
    border: none;
  }
  button:first-child {
    background-color: white;
  }
  .dropdown-item {
    &.active {
      background-color: $app-bg2;
    }
  }
  .custom-dropdown {
    position: absolute;
    top: 43px;
    z-index: 15;
  }

  // Secondary button.
  .dropdown {
    .search--btn {
      // @include app-nav-kit;
      :hover {
        @extend %app-nav-color;
      }
    }
    .btn-secondary {
      border: none;
      background-color: none;
    }
    .dropdown-menu {
      width: 100%;
    }
  }
}

</style>
