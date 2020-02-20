<template>
  <b-container fluid class="">
    <!-- Header -->
    <b-row>
      <b-col>
        <b-navbar class="app-header"
           toggleable="md">
          <b-navbar-brand href="#" class="pr-5">
            <b-img src="~@/assets/logo.png"
               alt="Knaxim Logo" />
            Knaxim.com
          </b-navbar-brand>
          <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
          <b-collapse id="nav-collapse" is-nav>

           <!-- Search & History -->
           <header-search />

           <!-- Settings Nav -->
           <header-settings @login="showLogin"/>
          </b-collapse>
        </b-navbar>
      </b-col>
      <hr class="w-100 m-0"/>
    </b-row>

    <b-row>
      <!-- Side Nav -->
      <b-col class="app-side pl-0 mr-2 min-max-150" cols="2">
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
      </b-col>

      <b-col class="overflow-auto">
        <!-- Sub Header -->
        <b-row>
          <b-col>
            <team-select v-if="isAuthenticated" class="teamselect"
              @team-selected="gotoTeam"/>
          </b-col>
        </b-row>

        <auth ref="auth" :passkey="resetkey"></auth>
        <b-row class="">
          <!-- Main Content -->
          <b-col class="p-0">
            <div class="app-content">
              <div v-if="!isAuthenticated">
                <div class="empty">
                  <h1>You aren't logged in!</h1>
                  <b-button @click="showAuth">
                    <h3>Login</h3>
                  </b-button>
                </div>
              </div>
              <router-view v-else/>
            </div>
          </b-col>

          <!-- Side View -->
          <router-view name="sideview" />
        </b-row>
      </b-col>
    </b-row>

    <!-- TODO: A more descriptive error Object vs error String, would
    be useful for designing better UI error segments. -->
    <error-control global v-on:error="makeToast($event)">
      <div></div>
    </error-control>
  </b-container>

</template>

<script>
// Header
import HeaderSearch from '@/components/header-search'
import HeaderSettings from '@/components/header-settings'

// Sub Header
import TeamSelect from '@/components/team-select'
import { mapGetters } from 'vuex'
import { GET_USER, LOAD_SERVER } from '@/store/actions.type'

// Side Nav
import NavBasic from '@/components/nav-basic'
import StorageInfo from '@/components/storage-info'
import HeaderSearchHistory from '@/components/header-search-history'

import Auth from '@/components/auth'
import ErrorControl from '@/components/error-control'

export default {
  name: 'App',
  data () {
    return {
      appInfoDisplay: null,
      context: 'My Cloud',
      auth: false
    }
  },
  created () {
    this.$store.dispatch(GET_USER).then(() => {
      this.$store.dispatch(LOAD_SERVER)
    }).catch(() => {
      this.showAuth()
    })
  },
  methods: {
    showAuth () {
      if (this.$route.name === 'reset') {
        this.$refs.auth.openReset()
      } else {
        this.$refs.auth.openLogin()
      }
    },
    // Used by ErrorControl to display specific errors: login, ...
    makeToast (msg, append = false) {
      this.$bvToast.toast(msg, {
        title: 'Error',
        autoHideDelay: 5000,
        appendToast: append,
        ...(msg === 'Please Login.' ? { to: '/login' } : '')
      })
    },

    gotoTeam (id) {
      if (id === this.currentUser.id) {
        this.$router.push({ name: 'home' })
      } else {
        this.$router.push(`/team/${id}`)
      }
    },

    showLogin () {
      this.$refs.auth.openLogin()
    }
  },

  computed: {
    resetkey () {
      return this.$route.params ? this.$route.params.passkey || '' : ''
    },
    ...mapGetters(['isAuthenticated', 'currentUser'])
  },

  components: {
    // Header
    HeaderSearch,
    HeaderSettings,

    // Sub Header
    TeamSelect,

    // Side Nav
    NavBasic,
    StorageInfo,
    HeaderSearchHistory,

    Auth,
    ErrorControl
  }
}
</script>

<style lang="scss">
/* TODO: Is this the correct place to style the BODY...background-color is being
overridden by something after this point. */
body {
  overflow-y: hidden;
  // background-color: #e5e5e5 !important;
}

#app {
  hr {
    border-top: 1px solid rgba(0, 0, 0, 0.4);
  }
  .col {
    // border: 1px dashed red;
  }
}

/* TODO: Use Flexbox/Bootstrap-vue/Bootstrap to grow the app to 100% */
.app-content {
  height: calc(100vh - 140px);
}

/* App Header */
.app-header {
  img {
    width: 50px;
    height: 50px;
  }
}

/* Sub Nav */
.app-subnav {
  height: 55px;
}
.teamselect {
  max-width: 18%;
}
.empty {
  text-align: center;
  margin-top: 10%;
  button {
    background-color: white;
    border-radius: 10px;
    border: 0px;
    width: 160px;
    height: 80px;
    color: rgb(46, 46, 46);
  }

  button:hover {
    background-color: rgb(150, 182, 252);
    color: rgb(46, 46, 46);
  }
}
</style>
