<!--
auth.vue: component for handling and juggling the long and registration modals

call this component's openLogin() method to open the login modal
-->
<template>
  <div>
    <login-modal :userFill="userFill" ref="login" id="auth-login"
      @register="pushReg" @login="loginSuccess" @close="loginClose"/>
    <registration-modal ref="reg" id="auth-register" @close="pushLogin" @register="regSuccess"/>
  </div>
</template>

<script>
import LoginModal from '@/components/modals/login-modal'
import RegistrationModal from '@/components/modals/registration-modal'
import { FILES_LIST } from '@/store/actions.type'

export default {
  name: 'auth',
  components: {
    LoginModal,
    RegistrationModal
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
    openReg () {
      this.$refs['reg'].show()
    },
    pushReg () {
      this.$router.push('/register')
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
    }
  },
  mounted () {
    if (this.$route.name === 'login') {
      this.openLogin()
    } else if (this.$route.name === 'register') {
      this.openReg()
    }
  },
  watch: {
    $route (to, from) {
      if (to.name === 'login') {
        this.openLogin()
      } else if (to.name === 'register') {
        this.openReg()
      }
    }
  }
}
</script>

<style>

</style>
