<template>
  <team-control v-if="!!activeGroup" #default="{ owner, members, addMember, removeMember }">
    <b-col class="member-list" cols="2">
      <b-list-group>
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
        <b-input name="team-name" v-model="name"/>
      </b-modal>
    </b-col>
  </team-control>
</template>

<script>
import TeamControl from '@/components/team-control'
import { mapGetters } from 'vuex'

export default {
  name: 'member-list',
  components: {
    TeamControl
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
  }
</style>
