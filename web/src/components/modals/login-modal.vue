<!--
login-modal: a window for logging in

events:
  'login': emitted upon successful login
  'register': emitted upon the 'Create Account' button being clicked
  'close': emitted upon the modal being closed

The methods show() and hide() are intended to be public methods for the parent.
-->
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
    <b-img src="@/assets/logo.png" alt="Knaxim Logo"/>
    <b-form @submit.prevent="login">
      <b-form-group :state="!fail" :invalid-feedback="feedback">
        <b-form-input v-model="username" placeholder="Username" autofocus ref="userField"/>
        <b-form-input v-model="password" placeholder="Password" type="password" ref="passField"/>
      </b-form-group>
      <div v-if="loading">
        <b-spinner class="m-3"/>
      </div>
      <b-form-group v-else>
        <b-button type="submit" class="shadow-sm">Login</b-button>
        <b-button @click="register" class="shadow-sm">Create Account</b-button>
        <b-button @click="request" class="shadow-sm">Forgot Password</b-button>
      </b-form-group>
    </b-form>
  </b-modal>
</template>

<script>
import { mapGetters } from 'vuex'
import { LOGIN } from '@/store/actions.type'

export default {
  name: 'login-modal',
  props: {
    id: {
      type: String,
      required: true
    },
    userFill: {
      type: String,
      required: false
    }
  },
  data () {
    return {
      username: '',
      password: '',
      fail: false
    }
  },
  computed: {
    feedback () {
      return 'Invalid username or password'
    },
    ...mapGetters({
      loading: 'authLoading'
    })
  },
  methods: {
    login () {
      this.$store.dispatch(LOGIN, { login: this.username, password: this.password }
      ).then((res) => {
        this.$emit('login')
        this.$refs['modal'].hide()
      }).catch((res) => {
        this.fail = true
      })
    },
    register () {
      this.$emit('register')
    },
    registerSuccess () {
      this.$refs.passField.focus()
    },
    show () {
      this.$refs['modal'].show()
    },
    hide () {
      this.$refs['modal'].hide()
    },
    onClose () {
      this.$emit('close')
    },
    request () {
      this.$emit('request')
    }
  },
  watch: {
    /*
    There have been issues with focusing the input fields during modal
    transitions. For one thing, $refs is undefined while the modal is hidden,
    and for another thing, focus may be lost right after changing input fields.
    Thus, a setTimeout() call is used with a wrapper for the focus.
    */
    userFill () {
      this.username = this.userFill
      setTimeout(this.registerSuccess, 200)
    }
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
