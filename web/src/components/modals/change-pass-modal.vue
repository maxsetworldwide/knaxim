<template>
  <b-modal
    :id="id"
    ref="modal"
    :no-close-on-backdrop="loading"
    :no-close-on-esc="loading"
    @hidden="onClose"
    centered
    hide-footer
    hide-header
    content-class="modal-style">
    <b-form @submit.prevent="changepass" v-if="isAuthenticated">
      <b-form-group>
        <b-form-input v-model="oldpass" placeholder="Current Password" autofocus ref="oldPass" type="password"/>
        <b-form-input v-model="newpass" placeholder="New Password" ref="newPass" type="password" :state="validPassword"/>
        <b-form-text>Password must have at least 6 characters, one capital, and one number.</b-form-text>
        <b-form-input v-model="newpassconfirm" placeholder="Confirm Password" ref="confirmPass" type="password" :state="matchPassword"/>
        <b-form-invalid-feedback>Password does not match.</b-form-invalid-feedback>
      </b-form-group>
      <div v-if="loading">
        <b-spinner class="m-3"/>
      </div>
      <b-form-group v-else :state="!fail" :invalid-feedback="feedback">
        <b-button type="submit" class="shadow-sm" :disabled="!validateForm">Change Password</b-button>
        <b-button @click="toLogin" class="shadow-sm">To Login</b-button>
      </b-form-group>
    </b-form>
    <b-form @submit.prevent="toLogin" v-else>
      <b-button type="submit" class="shadow">Please Login</b-button>
    </b-form>
  </b-modal>
</template>
<script>
import { mapGetters } from 'vuex'
import UserService from '@/service/user'

export default {
  name: 'change-pass-modal',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      oldpass: '',
      newpass: '',
      newpassconfirm: '',
      loading: false,
      fail: false
    }
  },
  computed: {
    feedback () {
      return 'invalid password'
    },
    validPassword () {
      if (this.newpass.length === 0) {
        return null
      }
      let digitRegex = /.*[0-9]+.*/
      let hasDigit = digitRegex.test(this.newpass)
      let hasCapital = this.newpass.toLowerCase() !== this.newpass
      return this.newpass.length >= 6 && hasDigit && hasCapital
    },
    matchPassword () {
      if (this.newpassconfirm.length === 0) {
        return null
      }
      return this.newpassconfirm === this.newpass
    },
    validateForm () {
      return this.matchPassword && this.validPassword && this.oldpass.length > 5
    },
    ...mapGetters(['isAuthenticated'])
  }
  methods: {
    changepass () {
      if (!this.validateForm) {
        return
      }
      this.loading = true
      UserService.changePassword({oldpass: this.oldpass, newpass: this.newpass}).then(res => {
        this.loading = false
        this.fail = false
        this.$emit('changed')
        this.hide()
      }, err => {
        this.loading = false
        this.fail = true
      })
      this.oldpass = ''
      this.newpass = ''
      this.newpassconfirm = ''
    },
    toLogin() {
      this.$router.push('/login')
    },
    onClose () {
      this.$emit('close')
    },
    hide () {
      this.$refs['modal'].hide()
    }
  }
}
</script>
