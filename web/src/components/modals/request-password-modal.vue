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
    <b-form @submit.prevent="send">
      <b-form-group>
        <b-form-input autofocus v-model="name" placeholder="Username"/>
      </b-form-group>
      <div v-if="loading">
        <b-spinner class="m-3"/>
      </div>
      <b-form-group v-else>
        <b-button type="submit" class="shadow-sm" :disabled="name.length < 6">Send Reset Email</b-button>
        <b-button @click="toLogin" class="shadow-sm">Back to Login</b-button>
      </b-form-group>
    </b-form>
  </b-modal>
</template>
<script>
import UserService from '@/service/user'

export default {
  name: 'request-password-modal',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      name: '',
      loading: false
    }
  },
  methods: {
    close () {
      this.$emit('close')
    },
    send () {
      if (this.name.length < 6) {
        return
      }
      this.loading = true
      UserService.requestReset({ name: this.name }).then(res => {
        this.toLogin()
        this.loading = false
      }).catch(() => {
        this.loading = false
      })
    },
    toLogin () {
      this.close()
    },
    show () {
      this.$refs['modal'].show()
    },
    hide () {
      this.name = ''
      this.$refs['modal'].hide()
    }
  }
}
</script>
