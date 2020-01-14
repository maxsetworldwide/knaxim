<template>
  <team-select v-if="isAuthenticated" class="teamselect" @team-selected="gotoTeam"/>
</template>

<script>
import teamSelect from '@/components/team-select'
import { mapGetters } from 'vuex'

export default {
  components: {
    teamSelect
  },
  name: 'app-subnav',
  props: {
    context: {
      default: 'My Cloud',
      type: String
    }
  },
  data () {
    return {
    }
  },
  computed: {
    ...mapGetters(['isAuthenticated', 'currentUser'])
  },
  methods: {
    gotoTeam (id) {
      if (id === this.currentUser.id) {
        this.$router.push({ name: 'home' })
      } else {
        this.$router.push(`/team/${id}`)
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.app-subnav {
  height: 55px;
}
.teamselect {
  max-width: 18%;
}
</style>
