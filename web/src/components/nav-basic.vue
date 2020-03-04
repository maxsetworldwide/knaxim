<template>
  <div class="nav-basic">
    <nav-basic-add class="p-4 px-4 pb-4 d-none d-md-flex min-max-150" />

    <b-nav vertical id="normal-nav" class="pb-4 min-max-150 d-none d-md-block">
      <b-nav-item class="mb-1" to="/" exact>
        <b-icon icon="inbox" class="icon" />
        <span>All</span>
      </b-nav-item>

      <b-nav-item class="mb-1" to="/list/owned" v-if="!groupMode">
        <b-icon icon="wallet" class="icon" />
        <span>{{ currentUser.name }}</span>
      </b-nav-item>

      <b-nav-item class="mb-1" to="/list/shared" v-if="!groupMode">
        <b-icon icon="people" class="icon" />
        <span>Shared</span>
      </b-nav-item>

      <b-nav-item class="mb-1" to="/list/recents">
        <b-icon icon="clock" class="icon" />
        <span>Recent</span>
      </b-nav-item>

      <b-nav-item class="mb-1" to="/list/favorites" v-if="!groupMode">
        <b-icon icon="heart" class="icon" />
        <span>Favorites</span>
      </b-nav-item>

      <b-nav-item class="mb-1" to="/list/trash" v-if="!groupMode">
        <b-icon icon="trash" class="icon" />
        <span>Trash</span>
      </b-nav-item>
    </b-nav>
    <b-navbar class="d-block d-md-none" tag="div">
      <b-navbar-nav fill align="center" class="small-flex-props">
        <nav-basic-add class="small-add" />
        <b-nav-text class="small-teamselect">
          <team-select class="teamselect" />
        </b-nav-text>
        <b-nav-dd ref="smallnav" id="small-nav-dd" :text="currentContext" right>
          <b-nav-item class="mb-1" to="/" @click="hide" exact>
            <b-icon icon="inbox" class="icon" />
            <span>All</span>
          </b-nav-item>

          <b-nav-item class="mb-1" to="/list/owned" @click="hide" v-if="!groupMode">
            <b-icon icon="wallet" class="icon" />
            <span>{{ currentUser.name }}</span>
          </b-nav-item>

          <b-nav-item class="mb-1" to="/list/shared" @click="hide" v-if="!groupMode">
            <b-icon icon="people" class="icon" />
            <span>Shared</span>
          </b-nav-item>

          <b-nav-item class="mb-1" to="/list/recents" @click="hide">
            <b-icon icon="clock" class="icon" />
            <span>Recent</span>
          </b-nav-item>

          <b-nav-item class="mb-1" to="/list/favorites" @click="hide" v-if="!groupMode">
            <b-icon icon="heart" class="icon" />
            <span>Favorites</span>
          </b-nav-item>

          <b-nav-item class="mb-1" to="/list/trash" @click="hide" v-if="!groupMode">
            <b-icon icon="trash" class="icon" />
            <span>Trash</span>
          </b-nav-item>
        </b-nav-dd>
      </b-navbar-nav>
    </b-navbar>
  </div>
</template>

<script>
import NavBasicAdd from '@/components/nav-basic-add.vue'
import TeamSelect from '@/components/team-select.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'nav-basic',
  props: {},
  components: {
    NavBasicAdd,
    TeamSelect
  },
  computed: {
    currentContext () {
      if (this.$route.name === 'home') {
        return '🗂 All'
      } else if (this.$route.name === 'filteredFiles') {
        switch (this.$route.params.src) {
          case 'recents':
            return '🗂 Recent'
          case 'favorites':
            return '🗂 Favorites'
          case 'shared':
            return '🗂 Shared'
          case 'owned':
            return '🗂 ' + this.currentUser.name
          case 'trash':
            return '🗂 Trash'
        }
      }
      return 'Select 🗂'
    },
    cloudtype () {
      if (this.activeGroup) {
        return 'Team'
      } else {
        return 'My'
      }
    },
    groupMode () {
      return !!this.activeGroup
    },
    ...mapGetters(['activeGroup', 'currentUser'])
  },
  methods: {
    hide () {
      this.$refs.smallnav.hide()
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
.nav-basic {
  #normal-nav {
    .nav-item {
      .nav-link {
        @extend %app-shadow-sm;
        @extend %app-nav-color;
        @extend %nav-right-round;

        padding-top: 4px;
        padding-bottom: 4px;
        width: 100%;
      }
      @include app-nav-kit;
    }
  }

  .small-flex-props {
    align-items: baseline;

  }

  .disabled {
    @extend %nav-disable;
  }

  .teamselect {
    width: 9em;
  }

  .small-add {
    padding-top: 5px;
    height: calc(100% - 5px);
    flex-grow: 1;
    width: 4em !important;
    margin-right: 1em;
  }

  .small-teamselect {
    flex-grow: 1;
    text-align: center;
  }

  .icon {
    @extend %nav-icon;
  }
}
</style>
