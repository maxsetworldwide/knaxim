<!--
Long term eventually this should be the root component for adjusting user profile options
-->
<template>
  <div>
    <change-pass-modal ref="changepass" id="change-password" @changed="home"/>
  </div>
</template>
<script>
import ChangePassModal from '@/components/modals/change-pass-modal'

export default {
  name: 'profile',
  components: {
    ChangePassModal
  },
  methods: {
    home () {
      this.$router.push('/')
    },
    showChangePass () {
      this.$refs['changepass'].show()
    }
  },
  mounted () {
    this.showChangePass()
  }
}
</script>
