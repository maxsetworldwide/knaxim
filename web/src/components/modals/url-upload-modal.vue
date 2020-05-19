<!--
url-upload-modal: window for submitting webpages

events:
  'upload': emitted upon successful upload
  'close': emitted upon any closure of the modal

global events:
  'url-upload': emitted along with 'upload'
-->
<template>
  <b-modal
    :id="id"
    @hidden="onClose"
    ref="modal"
    centered
    hide-footer
    title="Upload URL"
    :no-close-on-backdrop="loading"
    :no-close-on-esc="loading"
    content-class="modal-style"
  >
    <b-form @submit.prevent="upload">
      <b-form-input
        v-model="input"
        placeholder="Enter a web page URL"
        :state="validInput"
      />
      <b-form-text>
        Be sure to include http:// or https:// at the beginning of your URL!
      </b-form-text>
      <div v-if="loading">
        <b-spinner class="m-4" />
      </div>
      <div v-else>
        <b-button @click="upload" :disabled="!validInput" class="shadow-sm">
          Upload
        </b-button>
      </div>
    </b-form>
  </b-modal>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { CREATE_WEB_FILE } from '@/store/actions.type'
import { EventBus } from '@/plugins/utils'

export default {
  name: 'url-upload-modal',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      input: ''
    }
  },
  methods: {
    async upload () {
      if (!this.validInput) {
        return
      }
      this.loading = true
      await this.send({ url: this.input })
      this.$emit('upload')
      EventBus.$emit('url-upload')
      this.$refs['modal'].hide()
    },
    onClose () {
      this.$emit('close')
      this.input = ''
    },
    show () {
      this.$refs['modal'].show()
    },
    hide () {
      this.$refs['modal'].hide()
    },
    ...mapActions({
      send: CREATE_WEB_FILE
    })
  },
  computed: {
    validInput () {
      return this.input.length === 0
        ? null
        : this.input.indexOf('http') === 0 && this.input.indexOf('.') > 0
    },
    ...mapGetters({
      loading: 'fileLoading'
    })
  }
}
</script>

<style scoped lang="scss">
button {
  @extend %pill-buttons;
  width: flex;
  margin-right: 5px;
  margin-left: 5px;
  margin-top: 10px;
}

::v-deep .modal-style {
  @extend %modal-corners;
  text-align: center;
}
</style>
