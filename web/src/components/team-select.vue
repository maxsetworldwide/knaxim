<template>
  <b-form-select v-model="selected" :options="options"></b-form-select>
</template>
<script>
import { mapGetters } from 'vuex'
import { ACTIVATE_GROUP } from '@/store/mutations.type'
import { REFRESH_GROUPS } from '@/store/actions.type'

export default {
  name: 'team-select',
  data () {
    return {}
  },
  computed: {
    selected: {
      get () {
        if (this.currentGroup) {
          return this.currentGroup.id
        } else {
          return this.currentUser.id
        }
      },
      set (newVal) {
        this.onChange(newVal)
      }
    },
    options () {
      return [{ value: this.currentUser.id, text: this.currentUser.name }, ...this.availableGroups.map(group => {
        return {
          value: group.id,
          text: group.name
        }
      })]
    },
    ...mapGetters(['currentUser', 'availableGroups', 'activeGroup'])
  },
  methods: {
    onChange (id) {
      if (id === this.currentUser.id) {
        this.$store.commit(ACTIVATE_GROUP, { id: null })
      } else {
        this.$store.commit(ACTIVATE_GROUP, { id })
      }
    }
  },
  mounted () {
    this.$store.dispatch(REFRESH_GROUPS, {})
  }
}
</script>
