<template>
  <b-container fluid class="">
    <b-row>
      <b-col>
        <app-header />
      </b-col>
      <hr class="w-100 m-0"/>
    </b-row>

    <b-row>
      <b-col class="pl-0 mr-2 min-max-150" cols="2">
        <app-side />
      </b-col>

      <b-col class="overflow-auto">
        <b-row>
          <b-col>
            <app-subnav context="My Cloud" />
          </b-col>
        </b-row>

        <b-row class="">
          <b-col class="p-0">
            <div class="app-content">
              <!-- TODO: A more descriptive error Object vs error String, would
              be useful for designing better UI error segments. -->
              <error-control global v-on:error="makeToast($event)">
              </error-control>
              <router-view />
            </div>
          </b-col>
          <router-view name="sideview" />
        </b-row>
      </b-col>

    </b-row>
  </b-container>

</template>

<script>
import AppHeader from '@/components/app-header.vue'
import AppSubnav from '@/components/app-subnav.vue'
import AppSide from '@/components/app-side.vue'
import ErrorControl from '@/components/error-control'

export default {
  name: 'App',
  data () {
    return {
      appInfoDisplay: null
    }
  },
  methods: {
    makeToast (msg, append = false) {
      this.$bvToast.toast(msg, {
        title: 'Error',
        autoHideDelay: 5000,
        appendToast: append,
        ...(msg === 'Please Login.' ? { to: '/login' } : '')
      })
    }
  },
  computed: {
  },
  components: {
    AppHeader,
    AppSubnav,
    AppSide,
    ErrorControl
  },
  mounted () {
    // Capture all responses for login required.
    let vm = this
    this.axios.interceptors.response.use(function (config) {
      if (config.data && config.data.message === 'login') {
        vm.$root.$emit('app::set::error', 'Please Login.')
        // TODO: Throwing an error here blindly prevents duplicate requests
        // in some cases.  This may need to be handled better by the code
        // making requests, or perhaps a different mechanisim is needed.
        throw new Error('Login Required')
      }

      return config
    })
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

</style>
