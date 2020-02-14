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
import { FILES_LIST } from '@/store/actions.type'

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
      this.$router.push('/login')
      this.$refs['reg'].hide()
    },
    reqLogin () {
      this.$router.push('/login')
      this.$refs['request'].hide()
    },
    resLogin () {
      this.$router.push('/login')
      this.$refs['reset'].hide()
    },
    openReg () {
      this.$refs['reg'].show()
    },
    pushReg () {
      this.$router.push('/register')
      this.$refs['login'].hide()
    },
    pushRequest () {
      this.$router.push('/request')
      this.$refs['login'].hide()
    },
    regSuccess (username) {
      this.userFill = username
      this.openLogin()
    },
    loginSuccess () {
      this.$store.dispatch(FILES_LIST)
    },
    loginClose () {
      if (this.$route.name === 'login') {
        this.$router.push({ name: 'home' })
      }
    },
    openRequest () {
      this.$refs['request'].show()
    },
    openReset () {
      this.$refs['reset'].show()
    }
  },
  mounted () {
    if (this.$route.name === 'login') {
      this.openLogin()
    } else if (this.$route.name === 'register') {
      this.openReg()
    } else if (this.$route.name === 'request') {
      this.openRequest()
    } else if (this.$route.name === 'reset') {
      this.openReset()
    }
  },
  watch: {
    $route (to, from) {
      if (to.name === 'login') {
        this.openLogin()
      } else if (to.name === 'register') {
        this.openReg()
      } else if (to.name === 'request') {
        this.openRequest()
      } else if (to.name === 'reset') {
        this.openReset()
      }
    }
  }
}
</script>

<style>

</style>
