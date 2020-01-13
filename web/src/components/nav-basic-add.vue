<template>
  <b-dropdown
    id="nav-dropdown"
    text="Add"
    toggle-class="nav-basic-add--toggle"
    class="nav-basic-add w-100"
  >
    <!-- <b-dropdown-item>
      <svg>
        <use href="../assets/app.svg#folder-2" />
      </svg>
      <span>Folder</span>
    </b-dropdown-item> -->

    <team-control #default="{ createTeam }">
      <b-dropdown-item v-b-modal.at-modal>
        <svg>
          <use href="../assets/app.svg#group" />
        </svg>
        <span>Team</span>
        <b-modal id="at-modal" title="Create Team" @ok="createTeam({ name })">
          <b-input name="team-name" v-model="name" placeholder="Enter Team name"/>
        </b-modal>
      </b-dropdown-item>
    </team-control>

    <b-dropdown-divider />
    <b-dropdown-item v-b-modal.upload>
      <svg>
        <use href="../assets/app.svg#folder-2" />
      </svg>
      <span>File upload</span>
      <upload-modal id="upload"/>
    </b-dropdown-item>

    <!-- <b-dropdown-item>
      <svg>
        <use href="../assets/app.svg#folder-2" />
      </svg>
      <span>Folder upload</span>
    </b-dropdown-item> -->

    <b-dropdown-divider />
    <b-dropdown-item v-b-modal.batch-url-upload>
      <svg>
        <use href="../assets/app.svg#cloud" />
      </svg>
      <span>URL upload</span>
      <url-upload-modal id="batch-url-upload"/>
    </b-dropdown-item>
  </b-dropdown>
</template>

<script>
import UploadModal from '@/components/modals/upload-modal'
import UrlUploadModal from '@/components/modals/url-upload-modal'
import TeamControl from '@/components/team-control'

export default {
  name: 'nav-basic-add',
  components: {
    UploadModal,
    UrlUploadModal,
    TeamControl
  },
  props: {
  },
  data () {
    return {
      name: ''
    }
  },
  methods: {
    upload () {
      this.$bvModal.show('upload')
    },
    getUserOptions (users) {
      return users.map((member) => {
        return {
          text: member.name,
          value: member.id
        }
      })
    }
  }
}
</script>

<style lang="scss">

.nav-basic-add {
  .dropdown-menu {
    @extend %app-shadow-sm;
  }
  svg {
    height: 25px;
    width: 25px;
    margin-right: 15px;
  }
  .dropdown-item {
    padding-left: .5rem;
  }

  // Activating button.
  .btn-secondary:not(:disabled):not(.disabled):active {
    background-color: $app-clr1;
    color: $app-clr3;
    border: none;
  }

  // Active button.
  &.show > .btn-secondary.dropdown-toggle {
    @extend %app-shadow-sm;
    @extend %app-nav-color;
    background-color: $app-bg2;
    border: none;
  }

  .btn-secondary {
    @extend %nav-round;
    box-shadow: 0 0.1rem 0.25rem 3px rgba(102, 255, 51, 0.4);
    @extend %app-nav-color;
    border: none;

    &:hover {
      background-color: $app-bg4;
      border: none;
      color: $app-clr;
    }
  }
}

</style>
