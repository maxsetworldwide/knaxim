<template>
  <div class="nav-basic min-max-150">
    <nav-basic-add class="p-4 px-4 pb-4"/>

    <b-nav vertical class="pb-4">
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

      <!-- <b-nav-item class="mb-1"
          disabled to="/user/file/delete">
        <svg>
          <use href="../assets/app.svg#bin" />
        </svg>
        <span>Delete</span>
      </b-nav-item> -->
    </b-nav>

  </div>
</template>

<script>
import NavBasicAdd from '@/components/nav-basic-add.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'nav-basic',
  props: {
  },
  components: {
    NavBasicAdd
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

  .disabled {
    @extend %nav-disable;
  }
}
</style>
