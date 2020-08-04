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
  <team-control v-if="!!activeGroup" #default="{ owner, members, addMember, removeMember }">
    <b-col class="member-list" cols="2">
      <b-list-group  class="team-members">
        <b-list-group-item>
          <svg>
            <use href="../assets/app.svg#group" />
          </svg>
          <span> {{ activeGroup.name }}</span>
          <br>
          <span> Leader: {{ owner.name }} </span>
        </b-list-group-item>
        <b-list-group-item v-for="mem in members" :key="mem.id">
          {{ mem.name }}
          <b-button
           @click="removeMember(mem.id)"
           pill
           variant="outline-danger"
           size="sm"
           v-if="currentUser.id === owner.id"
          >remove</b-button>
        </b-list-group-item>
        <b-list-group-item v-if="currentUser.id === owner.id">
          <b-button v-b-modal.add-modal
          variant="primary" size="sm">Add</b-button>
        </b-list-group-item>
      </b-list-group>
      <b-modal id="add-modal" title="Add Member" @ok="addMember(name)">
        <b-input list="team-control-input" name="team-name" v-model="name"/>
        <owner-dropdown id="team-control-input" />
      </b-modal>
    </b-col>
  </team-control>
</template>

<script>
import TeamControl from '@/components/team-control'
import OwnerDropdown from '@/components/owner-dropdown'
import { mapGetters } from 'vuex'

export default {
  name: 'member-list',
  components: {
    TeamControl,
    OwnerDropdown
  },
  props: {
    selected: {
      type: Array,
      default: function () {
        return []
      }
    }
  },
  data () {
    return {
      processing: false,
      name: ''
    }
  },
  computed: {
    ...mapGetters(['activeGroup', 'currentUser'])
  },
  methods: {
  }
}
</script>

<style lang="scss">
  .member-list {
    svg {
      height: 50px;
      width: 50px;
    }

    .team-members {
      position: absolute;
      max-height: 100%;
      overflow-y: auto;
    }
  }
</style>
