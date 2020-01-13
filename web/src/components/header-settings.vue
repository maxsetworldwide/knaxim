<template>
 <!-- Settings Gear -->
 <b-navbar-nav class="header-settings">
   <b-nav-item-dropdown class="ml-auto pl-4" no-caret>
     <template #button-content>
       <svg>
         <use href="../assets/app.svg#settings" />
       </svg>
     </template>

     <b-dropdown-item @click="login">Login</b-dropdown-item>
     <b-dropdown-item @click="logout">Logout</b-dropdown-item>
   </b-nav-item-dropdown>
 </b-navbar-nav>
</template>

<script>
import { LOGOUT } from '@/store/actions.type'

export default {
  name: 'header-settings',
  data () {
    return {
    }
  },
  props: {
  },
  components: {
  },
  methods: {
    login () {
      this.$router.push('/login')
    },
    logout () {
      this.$store.dispatch(LOGOUT)
    }
  }
}
</script>

<style lang="scss">
.header-settings {
  svg {
    width: 35px;
    height: 35px;
  }
}
</style>
