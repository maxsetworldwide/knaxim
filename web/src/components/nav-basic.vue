<template>
  <div class="nav-basic">
    <nav-basic-add class="p-4 px-4 pb-4 d-none d-md-flex"/>

    <b-nav vertical id="normal-nav" class="pb-4 min-max-150 d-none d-md-block">
      <b-nav-item class="mb-1"
          to="/" exact>
        <svg>
          <use href="../assets/app.svg#cloud" />
        </svg>
        <span>{{ cloudtype }} Cloud</span>
      </b-nav-item>

      <b-nav-item class="mb-1"
          to="/list/owned"
          v-if="!groupMode">
        <svg>
          <use href="../assets/app.svg#files" />
        </svg>
        <span>Owned</span>
      </b-nav-item>

      <b-nav-item class="mb-1"
          to="/list/shared"
          v-if="!groupMode">
        <svg>
          <use href="../assets/app.svg#transfer" />
        </svg>
        <span>Shared</span>
      </b-nav-item>

      <b-nav-item class="mb-1"
          to="/list/recents">
        <svg>
          <use href="../assets/app.svg#clock" />
        </svg>
        <span>Recent</span>
      </b-nav-item>

      <b-nav-item class="mb-1"
          to="/list/favorites"
          v-if="!groupMode">
        <svg>
          <use href="../assets/app.svg#star" />
        </svg>
        <span>Favorites</span>
      </b-nav-item>

      <b-nav-item class="mb-1"
          to="/list/trash"
          v-if=!groupMode>
        <svg>
          <use href="../assets/app.svg#bin" />
        </svg>
        <span>Trash</span>
      </b-nav-item>
    </b-nav>
    <b-navbar class="d-block d-md-none" tag="div">
      <b-navbar-nav fill align="center">
        <nav-basic-add left class="small-add"/>
        <b-nav-text>
          <team-select class="teamselect"/>
        </b-nav-text>
        <b-nav-dd
          id="small-nav-dd"
          text=" 🗂 Context"
          right
          >
          <b-nav-item class="mb-1"
              to="/" exact>
            <svg>
              <use href="../assets/app.svg#cloud" />
            </svg>
            <span>{{ cloudtype }} Cloud</span>
          </b-nav-item>

          <b-nav-item class="mb-1"
              to="/list/owned"
              v-if="!groupMode">
            <svg>
              <use href="../assets/app.svg#files" />
            </svg>
            <span>Owned</span>
          </b-nav-item>

          <b-nav-item class="mb-1"
              to="/list/shared"
              v-if="!groupMode">
            <svg>
              <use href="../assets/app.svg#transfer" />
            </svg>
            <span>Shared</span>
          </b-nav-item>

          <b-nav-item class="mb-1"
              to="/list/recents">
            <svg>
              <use href="../assets/app.svg#clock" />
            </svg>
            <span>Recent</span>
          </b-nav-item>

          <b-nav-item class="mb-1"
              to="/list/favorites"
              v-if="!groupMode">
            <svg>
              <use href="../assets/app.svg#star" />
            </svg>
            <span>Favorites</span>
          </b-nav-item>

          <b-nav-item class="mb-1"
              to="/list/trash"
              v-if=!groupMode>
            <svg>
              <use href="../assets/app.svg#bin" />
            </svg>
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
  props: {
  },
  components: {
    NavBasicAdd,
    TeamSelect
  },
  computed: {
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
    ...mapGetters(['activeGroup'])
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
        svg {
          height: 25px;
          width: 25px;
          margin-right: 15px;
        }
      }
      @include app-nav-kit;
    }
  }

  .disabled {
    @extend %nav-disable;
  }

  svg {
    height: 25px;
    width: 25px;
    margin-right: 15px;
  }

  .teamselect {
    width: 9em;
  }
  .small-add {
    padding-top: 5px;
    height: calc(100% - 5px);
  }
}
</style>
