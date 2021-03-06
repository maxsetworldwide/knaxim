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
<script>
import { mapActions, mapGetters } from 'vuex'
import { CREATE_GROUP, LOAD_OWNER, ADD_MEMBER, REMOVE_MEMBER, LOOKUP_OWNER } from '@/store/actions.type'

export default {
  name: 'team-control',
  components: {},
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
      processing: false
    }
  },
  computed: {
    members () {
      if (this.activeGroup) {
        return this.activeGroup.members.map(id => {
          this[LOAD_OWNER]({ id })
          return {
            id,
            name: this.ownerNames[id] || 'loading...'
          }
        })
      }
      return []
    },
    owner () {
      if (this.activeGroup) {
        this[LOAD_OWNER]({ id: this.activeGroup.owner })
        return {
          id: this.activeGroup.owner || '',
          name: this.ownerNames[this.activeGroup.owner] || 'loading...'
        }
      }
      return {
        id: '',
        name: ''
      }
    },
    loading () {
      return this.ownerLoading || this.groupLoading
    },
    ...mapGetters(['activeGroup', 'ownerNames', 'ownerLoading', 'groupLoading'])
  },
  methods: {
    ...mapActions([CREATE_GROUP, ADD_MEMBER, REMOVE_MEMBER, LOOKUP_OWNER, LOAD_OWNER]),
    removeGroup () {},
    addTeamMember (targetname) {
      this[LOOKUP_OWNER]({ name: targetname })
        .then(id => {
          this[ADD_MEMBER]({ newMember: id })
        })
        .catch(() => {
          this.$bvToast.toast(`${targetname} does not exist`, {
            title: 'Unable to add member',
            autoHideDelay: 5000,
            appendToast: false
          })
        })
    },
    removeTeamMember (target) {
      this[REMOVE_MEMBER]({ newMember: target })
    }
  },

  render () {
    return this.$scopedSlots.default({
      owner: this.owner,
      members: this.members,
      createTeam: this[CREATE_GROUP],
      removeTeam: this.removeGroup,
      addMember: this.addTeamMember,
      removeMember: this.removeTeamMember
    })
  }
}
</script>
