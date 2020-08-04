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
<!--
registration-modal: a window for creating an account

events:
  'register': emitted upon successful registration, passing the username.
  'cancel': emitted upon closing with no successful registration.
  'close': emitted upon any closure of the modal.
-->
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
    <b-form @submit.prevent="register">
      <b-form-group>
        <b-form-input autofocus v-model="email" placeholder="Email" :state="validEmail"/>
        <b-form-invalid-feedback force-show v-if="fail">invalid email, or username already in use</b-form-invalid-feedback>
        <b-form-group>
          <b-form-input v-model="username" placeholder="Username" :state="validUser"/>
          <b-form-invalid-feedback>
            {{ nameTaken ? 'Name taken!' : 'Username must be between 6 and 12 characters long.' }}
          </b-form-invalid-feedback>
        </b-form-group>
        <b-form-group>
          <b-form-input v-model="password" placeholder="Password" type="password" :state="validPassword"/>
          <b-form-text>
            Password must have at least 6 characters, one capital, and one number.
          </b-form-text>
          <b-form-input v-model="passConf" placeholder="Confirm Password" type="password" :state="matchPasswords"/>
          <b-form-invalid-feedback>Passwords must match</b-form-invalid-feedback>
        </b-form-group>
        <b-form-group>
          <b-form-checkbox
            id="checkbox-terms"
            v-model="acceptedTerms"
            class="accept-terms"
            value="accepted"
            unchecked-value="not_accepted"
          >I understand and agree to the <a href="/APITermsOfService.pdf">Terms of Service</a> and <a href="/CustomerTermsOfService.pdf">Customer Terms of Service</a></b-form-checkbox>
        </b-form-group>
      </b-form-group>
      <div v-if="loading">
        <b-spinner class="m-3"/>
      </div>
      <b-form-group v-else>
        <b-button type="submit" class="shadow-sm" :disabled="!validateForm">Create Account</b-button>
        <b-button @click="hide" class="shadow-sm">Back to Login</b-button>
      </b-form-group>
    </b-form>
  </b-modal>
</template>

<script>
import { REGISTER } from '@/store/actions.type'

export default {
  name: 'registration-modal',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      email: '',
      username: '',
      password: '',
      passConf: '',
      loading: false,
      takenNames: [],
      success: false,
      fail: false,
      acceptedTerms: 'not_accepted'
    }
  },
  methods: {
    register () {
      if (!this.validateForm) {
        return
      }
      this.loading = true
      this.$store.dispatch(REGISTER, { email: this.email, login: this.username, password: this.password }
      ).then((res) => {
        this.loading = false
        this.success = true
        this.$emit('register', this.username)
        this.hide()
      }).catch((res) => {
        this.loading = false
        if (res.message === 'Name Already Taken') {
          this.takenNames.push(this.username)
        }
        this.fail = true
      })
    },
    show () {
      this.$refs['modal'].show()
    },
    hide () {
      this.password = ''
      this.passConf = ''
      this.$refs['modal'].hide()
    },
    close () {
      if (!this.success) {
        this.$emit('cancel')
      }
      this.$emit('close')
    }
  },
  computed: {
    validateForm () {
      return this.validEmail && this.validUser && this.validPassword && this.matchPasswords && this.acceptedTerms === 'accepted'
    },
    validEmail () {
      if (this.email.length === 0) {
        return null
      }
      const re = /.+@.+/
      return this.email.match(re) !== null
    },
    nameTaken () {
      return this.takenNames.indexOf(this.username) > -1
    },
    validUser () {
      if (this.username.length === 0) {
        return null
      }
      return this.username.length >= 6 && this.username.length <= 12 && !this.nameTaken
    },
    validPassword () {
      if (this.password.length === 0) {
        return null
      }
      let digitRegex = /.*[0-9]+.*/
      let hasDigit = digitRegex.test(this.password)
      let hasCapital = this.password.toLowerCase() !== this.password
      return this.password.length >= 6 && hasDigit && hasCapital
    },
    matchPasswords () {
      if (this.passConf.length === 0) {
        return null
      }
      return this.password === this.passConf
    }
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

.accept-terms {
  width: 80%;
  margin-left: 10%;
  margin-right: 10%;
  margin-top: 10px;
  margin-bottom: 10px;
}

</style>
