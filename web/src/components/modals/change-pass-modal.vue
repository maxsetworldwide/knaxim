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
    <b-img src="@/assets/CloudEdison.png" alt="Cloud Edison"/>
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
        <b-button @click="onClose" class="shadow-sm">Close</b-button>
      </b-form-group>
    </b-form>
    <b-form @submit.prevent="toLogin" v-else>
      <b-button type="submit" class="shadow">Please Login</b-button>
    </b-form>
  </b-modal>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import { CHANGE_PASSWORD } from '@/store/actions.type'

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
    ...mapGetters(['isAuthenticated']),
    ...mapGetters({
      loading: 'authLoading'
    })
  },
  methods: {
    changepass () {
      if (!this.validateForm) {
        return
      }
      this.send({ oldpass: this.oldpass, newpass: this.newpass }).then(() => {
        this.fail = false
        this.$emit('changed')
        this.hide()
      }, () => {
        this.fail = true
      }).finally(() => {
        this.oldpass = ''
        this.newpass = ''
        this.newpassconfirm = ''
      })
    },
    show () {
      this.$refs['modal'].show()
    },
    onClose () {
      this.$emit('close')
    },
    hide () {
      this.$refs['modal'].hide()
    },
    ...mapActions({ send: CHANGE_PASSWORD })
  }

}
</script>

<style scoped lang="scss">

img {
  width: 50%;
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
