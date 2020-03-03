<template>
  <b-modal
  :id="id"
  ref="modal"
  :no-close-on-backdrop="loading"
  :no-close-on-esc="loading"
  @hidden="close"
  centered
  hide-footer
  hide-header
  content-class="modal-style">
    <b-img src="@/assets/CloudEdison.png" alt="Cloud Edison"/>
    <b-form @submit.prevent="sendReset">
      <b-form-group>
        <b-input v-model="newpass" placeholder="Password" type="password" :state="validPassword" autofocus/>
        <b-form-text>
          Password must have at least 6 characters, one capital, and one number.
        </b-form-text>
        <b-form-input v-model="passConf" placeholder="Confirm Password" type="password" :state="matchPasswords"/>
        <b-form-invalid-feedback>Passwords must match</b-form-invalid-feedback>
      </b-form-group>
      <div v-if="loading">
        <b-spinner class="m-3"/>
      </div>
      <b-form-group v-else>
        <b-button type="submit" class="shadow-sm" :disabled="!validateForm">Update Password</b-button>
        <b-button @click="close" class="shadow-sm">Back to Login</b-button>
      </b-form-group>
    </b-form>
  </b-modal>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import { RESET_PASSWORD } from '@/store/actions.type'

export default {
  name: 'reset-password-modal',
  props: {
    id: {
      type: String,
      required: true
    },
    passkey: {
      type: String,
      required: true,
      default: ''
    }
  },
  data () {
    return {
      newpass: '',
      passConf: ''
    }
  },
  methods: {
    show () {
      this.$refs['modal'].show()
    },
    hide () {
      this.newpass = ''
      this.passConf = ''
      this.$refs['modal'].hide()
    },
    close () {
      this.$emit('close')
    },
    sendReset () {
      if (!this.validateForm) {
        return
      }
      this.loading = true
      this.reset({ passkey: this.passkey, newpass: this.newpass }).then(() => {
        this.close()
      })
    },
    ...mapActions({
      reset: RESET_PASSWORD
    })
  },
  computed: {
    validateForm () {
      return this.validPassword && this.matchPasswords
    },
    validPassword () {
      if (this.newpass.length === 0) {
        return null
      }
      let digitRegex = /.*[0-9].*/
      let hasDigit = digitRegex.test(this.newpass)
      let hasCapital = this.newpass.toLowerCase() !== this.newpass
      return this.newpass.length >= 6 && hasDigit && hasCapital
    },
    matchPasswords () {
      if (this.passConf.length === 0) {
        return null
      }
      return this.newpass === this.passConf
    },
    ...mapGetters({
      loading: 'authLoading'
    })
  }
}
</script>

<style scoped lang="scss">

img {
  width: 30%;
}

input {
  margin-top: 10px;
  margin-bottom: 10px;
  width: 80%;
  display: inline-block;
}

button {
  @extend %pill-buttons;
  width: flex;
  margin-right: 5px;
  margin-left: 5px;
}

::v-deep .modal-style {
  @extend %modal-corners;
  text-align: center;
}

</style>
