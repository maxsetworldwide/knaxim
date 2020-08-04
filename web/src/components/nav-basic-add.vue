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
  <b-dropdown
    id="nav-dropdown"
    text="Add"
    toggle-class="nav-basic-add--toggle"
    class="nav-basic-add w-100"
  >
    <team-control #default="{ createTeam }">
      <b-dropdown-item @click="showTeamModal">
        <svg>
          <use href="../assets/app.svg#group" />
        </svg>
        <span>Team</span>
        <b-modal
          id="at-modal"
          ref="at-modal"
          title="Create Team"
          @ok="createTeam({ name })"
        >
          <b-input
            name="team-name"
            v-model="name"
            placeholder="Enter Team name"
          />
        </b-modal>
      </b-dropdown-item>
    </team-control>

    <b-dropdown-divider />
    <b-dropdown-item @click="showUploadModal">
      <svg>
        <use href="../assets/app.svg#folder-2" />
      </svg>
      <span>File upload</span>
      <upload-modal id="upload-modal" ref="upload-modal" />
    </b-dropdown-item>

    <!-- <b-dropdown-item>
      <svg>
        <use href="../assets/app.svg#folder-2" />
      </svg>
      <span>Folder upload</span>
    </b-dropdown-item> -->

    <b-dropdown-divider />
    <b-dropdown-item @click="showURLModal">
      <b-icon icon="window" />
      <span>URL upload</span>
      <url-upload-modal id="url-modal" ref="url-modal" />
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
  props: {},
  data () {
    return {
      name: ''
    }
  },
  methods: {
    getUserOptions (users) {
      return users.map(member => {
        return {
          text: member.name,
          value: member.id
        }
      })
    },
    showTeamModal () {
      this.$refs['at-modal'].show()
    },
    showUploadModal () {
      this.$refs['upload-modal'].show()
    },
    showURLModal () {
      this.$refs['url-modal'].show()
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
    fill: $app-icon;
  }
  .dropdown-item {
    padding-left: 0.5rem;
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
