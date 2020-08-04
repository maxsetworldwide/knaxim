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
auth.vue: component for handling and juggling the long and registration modals

call this component's openLogin() method to open the login modal
-->
<template>
  <div>
    <login-modal :userFill="userFill" ref="login" id="auth-login"
      @register="pushReg" @login="loginSuccess" @request="pushRequest" @close="loginClose"/>
    <registration-modal ref="reg" id="auth-register" @close="pushLogin" @register="regSuccess"/>
    <request-password-modal id="auth-request" ref="request" @close="reqLogin"/>
    <reset-password-modal id="auth-reset" :passkey="passkey" ref="reset" @close="resLogin"/>
  </div>
</template>

<script>
import LoginModal from '@/components/modals/login-modal'
import RegistrationModal from '@/components/modals/registration-modal'
import RequestPasswordModal from '@/components/modals/request-password-modal'
import ResetPasswordModal from '@/components/modals/reset-password-modal'

export default {
  name: 'auth',
  components: {
    LoginModal,
    RegistrationModal,
    RequestPasswordModal,
    ResetPasswordModal
  },
  props: {
    passkey: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      userFill: ''
    }
  },
  methods: {
    openLogin () {
      this.$refs['login'].show()
    },
    pushLogin () {
      this.openLogin()
      this.$refs['reg'].hide()
    },
    reqLogin () {
      this.openLogin()
      this.$refs['request'].hide()
    },
    resLogin () {
      this.$refs['reset'].hide()
      this.$refs.login.show()
    },
    openReg () {
      this.$refs['reg'].show()
    },
    pushReg () {
      this.openReg()
      this.$refs['login'].hide()
    },
    pushRequest () {
      this.$refs.request.show()
      this.$refs['login'].hide()
    },
    regSuccess (username) {
      this.userFill = username
      this.openLogin()
    },
    loginSuccess () {
      this.loginClose()
      // this.$store.dispatch(FILES_LIST)
    },
    loginClose () {
      this.$refs['login'].hide()
      // if (this.$route.name === 'login') {
      //   this.$router.push({ name: 'home' })
      // }
    },
    openRequest () {
      this.$refs['request'].show()
    },
    openReset () {
      this.$refs['reset'].show()
    }
  },
  mounted () {
    // if (this.$route.name === 'login') {
    //   this.openLogin()
    // } else if (this.$route.name === 'register') {
    //   this.openReg()
    // }
  },
  watch: {
    $route (to, from) {
      // if (to.name === 'login') {
      //   this.openLogin()
      // } else if (to.name === 'register') {
      //   this.openReg()
      // }
    }
  }
}
</script>

<style>

</style>
