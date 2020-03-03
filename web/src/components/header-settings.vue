<template>
 <!-- Settings Gear -->
 <b-navbar-nav class="header-settings">
   <b-nav-item-dropdown class="dropleft" no-caret>
     <template #button-content>
       <svg>
         <use href="../assets/app.svg#settings" />
       </svg>
     </template>
     <b-dropdown-item @click="goHome" class="d-block d-md-none">Home</b-dropdown-item>
     <b-dropdown-item @click="login">Login</b-dropdown-item>
     <b-dropdown-item @click="changePass">Change Password</b-dropdown-item>
     <b-dropdown-item @click="logout">Logout</b-dropdown-item>
   </b-nav-item-dropdown>
 </b-navbar-nav>
</template>

<script>
import { LOGOUT } from '@/store/actions.type'
import { mapGetters } from 'vuex'

export default {
  name: 'header-settings',
  data () {
    return {
    }
  },
  props: {
  },
  components: {
    ...mapGetters(['isAuthenticated'])
  },
  methods: {
    login () {
      this.$emit('login')
    },
    changePass () {
      this.$router.push('/profile/newpassword')
    },
    goHome () {
      this.$router.push('/')
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
