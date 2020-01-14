<template>
  <team-control v-if="!!activeGroup" #default="{ members, addMember, removeMember }">
    <b-col class="member-list" cols="2">
      <b-list-group>
        <b-list-group-item>
          <svg>
            <use href="../assets/app.svg#group" />
          </svg>
          <span> {{ activeGroup.name }}</span>
        </b-list-group-item>
        <b-list-group-item v-for="mem in members" :key="mem.id">
          {{ mem.name }}
          <button @click="removeMember(mem.id)">remove</button>
        </b-list-group-item>
        <b-list-group-item>
          <button v-b-modal.add-modal>Add</button>
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
    ...mapGetters(['activeGroup'])
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
