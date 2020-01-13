<template>
  <div class="app-side">
    <b-row>
      <b-col>
        <nav-basic />
        <hr />
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <storage-info />
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <header-search-history />
      </b-col>
    </b-row>
  </div>
</template>

<script>
import NavBasic from '@/components/nav-basic.vue'
import StorageInfo from '@/components/storage-info.vue'
import HeaderSearchHistory from '@/components/header-search-history.vue'

export default {
  name: 'app-side',
  props: {
  },
  components: {
    NavBasic,
    StorageInfo,
    HeaderSearchHistory
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
.app-side {
}

</style>
