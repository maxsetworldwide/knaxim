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
    <b-form @submit.prevent="sendReset">
      <b-form-group>
        <b-input v-model="newpass" placeholder="Password" type="password" :state="validPassword"/>
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
import UserService from '@/service/user'

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
      passConf: '',
      loading: false
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
      UserService.resetPass({ passkey: this.passkey, newpass: this.newpass }).then(res => {
        this.loading = false
        this.close()
      }).catch(() => {
        this.loading = false
      })
    }
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
    }
  }
}
</script>
