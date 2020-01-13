<script>
import { mapActions, mapGetters } from 'vuex'
import { CREATE_GROUP } from '@/store/actions.type'
import GroupService from '@/service/group'
import UserService from '@/service/user'

export default {
  name: 'team-control',
  components: {},
  created () {
    this.loadMembers()
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
      members: [],
      loading: false,
      processing: false
    }
  },
  computed: {
    ...mapGetters(['activeGroup'])
  },
  methods: {
    ...mapActions([CREATE_GROUP]),
    removeGroup () {
    },
    addMember (targetname) {
      UserService.lookup({ name: targetname }).then(res => res.data).then(data => {
        if (data.id) {
          GroupService.add({ gid: this.activeGroup.id, target: data.id }).then(() => this.loadMembers())
        }
      })
      GroupService.lookup({ name: targetname }).then(res => res.data).then(data => {
        if (data.id) {
          GroupService.add({ gid: this.activeGroup.id, target: data.id }).then(() => this.loadMembers())
        }
      })
    },
    removeMember (target) {
      GroupService.remove({ gid: this.activeGroup.id, target }).then(() => this.loadMembers())
    },
    loadMembers () {
      if (this.activeGroup) {
        this.loading = true
        this.members = []
        let newmembers = this.members
        GroupService.info({ gid: this.activeGroup.id }).then(res => res.data)
          .then(data => {
            let memberGettingProms = []
            if (data.members) {
              data.members.forEach(memberid => {
                memberGettingProms.push(new Promise((resolve, reject) => {
                  let gp = GroupService.info({ gid: memberid }).then(res => res.data).then(data => {
                    if (data && data.id && data.name) {
                      newmembers.push(data)
                      resolve(true)
                    }
                  })
                  let up = UserService.info({ id: memberid }).then(res => res.data).then(data => {
                    if (data && data.id && data.name) {
                      newmembers.push(data)
                      resolve(true)
                    }
                  })
                  let inversegp = new Promise((resolve) => {
                    gp.catch(() => { resolve(false) })
                  })
                  let inverseup = new Promise((resolve) => {
                    up.catch(() => { resolve(false) })
                  })
                  Promise.all([inversegp, inverseup]).then(() => reject(new Error(`invalid id: ${memberid}`)))
                }))
              })
            }
            Promise.all(memberGettingProms).then(() => { this.loading = false; return true }).catch(() => {
              this.loading = false
            })
          })
      }
    }
  },
  watch: {
    activeGroup (n, o) {
      if (n !== o) {
        this.loadMembers()
      }
    }
  },

  render () {
    return this.$scopedSlots.default({
      members: this.members,
      createTeam: this[CREATE_GROUP],
      removeTeam: this.removeGroup,
      addMember: this.addMember,
      removeMember: this.removeMember
    })
  }
}
</script>
