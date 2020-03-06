<template>
  <b-container fluid class="">
    <!-- Header -->
    <b-row>
      <b-col>
        <auth ref="auth" :passkey="resetkey"></auth>
        <b-navbar class="app-header">
          <b-navbar-brand href="/" class="d-none d-md-block pr-5">
            <b-img src="~@/assets/CloudEdison.png" />
          </b-navbar-brand>

          <!-- Search & History -->
          <header-search />

          <!-- Settings Nav -->
          <header-settings @login="showLogin" />
        </b-navbar>
      </b-col>
      <hr class="w-100 m-0" />
    </b-row>

    <b-row>
      <!-- Side Nav -->
      <b-col v-if="isAuthenticated" class="flex-grow-0 d-none d-md-block">
        <b-row>
          <b-col>
            <b-row>
              <b-col>
                <nav-basic @team-selected="gotoTeam" />
              </b-col>
            </b-row>

            <b-row>
              <b-col>
                <storage-info class=""/>
              </b-col>
            </b-row>

            <b-row>
              <b-col>
                <header-search-history class="d-none d-md-block"/>
              </b-col>
            </b-row>
          </b-col>
        </b-row>
      </b-col>

      <b-col v-if="isAuthenticated" class="overflow-auto">
        <!-- Sub Header -->
        <b-row class="d-none d-md-flex">
          <b-col>
            <team-select
              class="teamselect"
              @team-selected="gotoTeam"
            />
          </b-col>
        </b-row>

        <b-row class="d-flex d-md-none">
          <b-col>
            <nav-basic @team-selected="gotoTeam" />
          </b-col>
        </b-row>

        <b-row class="app-content">
          <!-- Main Content -->
          <b-col class="p-0">
            <router-view />
          </b-col>

          <!-- Side View -->
          <router-view name="sideview" />
        </b-row>
      </b-col>

      <b-col v-else-if="loading" class="empty">
        <h1>Cloud Edison is Loading.</h1>
        <h1>Please Wait...</h1>
        <b-spinner style="width: 3em; height: 3em" />
      </b-col>

      <b-col v-else class="empty">
        <h1>You aren't logged in!</h1>
        <b-button @click="showAuth">
          <h3>Login</h3>
        </b-button>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
// Header
import HeaderSearch from '@/components/header-search'
import HeaderSettings from '@/components/header-settings'

// Sub Header
import TeamSelect from '@/components/team-select'
import { mapGetters, mapActions } from 'vuex'
import { GET_USER, LOAD_SERVER, ERROR_LOOP } from '@/store/actions.type'

// Side Nav
import NavBasic from '@/components/nav-basic'
import StorageInfo from '@/components/storage-info'
import HeaderSearchHistory from '@/components/header-search-history'

import Auth from '@/components/auth'

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
    this.$store
      .dispatch(GET_USER)
      .then(() => {
        this.$store.dispatch(LOAD_SERVER)
      })
      .catch(() => {
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
    makeToast (msg, title = 'Error', appendToast = false) {
      this.$bvToast.toast(msg, {
        title,
        noAutoHide: !!process.env.VUE_APP_DEBUG,
        autoHideDelay: 5000,
        appendToast
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
    },
    ...mapActions({
      handleErrors: ERROR_LOOP
    })
  },

  computed: {
    resetkey () {
      return this.$route.params ? this.$route.params.passkey || '' : ''
    },
    ...mapGetters([
      'isAuthenticated',
      'loading',
      'currentUser',
      'availableErrors'
    ])
  },
  watch: {
    availableErrors (newErrors) {
      if (newErrors) {
        if (process.env.VUE_APP_DEBUG) {
          this.handleErrors(e => this.makeToast(e.message, e.name || 'Error'))
        } else {
          this.handleErrors(() => {}) // Production drop errors
        }
      }
    }
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

    Auth
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
    width: 300px;
    height: 75px;
  }
}

/* Sub Nav */
.app-subnav {
  height: 55px;
}
.teamselect {
  max-width: 12em;
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
