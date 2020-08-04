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
