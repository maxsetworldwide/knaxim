<template>
 <b-navbar class="app-header"
    toggleable="md">
   <b-navbar-brand href="#" class="pr-5">
     <b-img src="~../assets/logo.png"
        alt="Knaxim Logo" />
     Knaxim.com
   </b-navbar-brand>

   <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
   <b-collapse id="nav-collapse" is-nav>

    <!-- Search & History -->
    <header-search />

    <!-- Settings Nav -->
    <header-settings />
   </b-collapse>
 </b-navbar>
</template>

<script>
import HeaderSearch from '@/components/header-search.vue'
import HeaderSettings from '@/components/header-settings.vue'

export default {
  name: 'app-header',
  data () {
    return {
    }
  },
  props: {
  },
  components: {
    HeaderSearch,
    HeaderSettings
  }
}
</script>

<style lang="scss">
.app-header {
  img {
    width: 50px;
    height: 50px;
  }
}
</style>
