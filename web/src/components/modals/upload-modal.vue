<!--
upload-modal: window for uploading files

events:
  'upload': emitted upon successful upload
  'close:': emitted upon any closure of the modal

global events:
  'file-upload': emitted along with 'upload'
-->
<template>
  <b-modal
    :id="id"
    ref="modal"
    @hidden="onClose"
    centered
    hide-footer
    hide-header
    :no-close-on-backdrop="loading"
    :no-close-on-esc="loading"
    content-class="modal-style"
  >
    <b-form-file
      v-model="files"
      v-bind:class="{ 'border-blue-shadow shadow-browse': files.length === 0 }"
      multiple
    >
      <template #file-name="{ names }">
        <b-badge>{{ names[0] }}</b-badge>
        <b-badge v-if="names.length > 1" class="ml-1">
          + {{ names.length - 1 }} More files
        </b-badge>
      </template>
    </b-form-file>

    <div v-if="loading">
      <b-spinner class="m-4" />
    </div>
    <div v-else>
      <b-button
        @click="upload"
        v-if="files.length > 0"
        :disabled="files.length === 0"
        class="border-blue-shadow"
        variant="primary"
      >
        Upload
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import { CREATE_FILE } from '@/store/actions.type'
import { EventBus } from '@/plugins/utils'

export default {
  name: 'upload-modal',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      files: [],
      loading: false
    }
  },
  methods: {
    upload () {
      this.loading = true
      let proms = []
      for (let i = 0; i < this.files.length; i++) {
        proms.push(this.$store.dispatch(CREATE_FILE, { file: this.files[i] }))
      }
      Promise.all(proms)
        .then(() => {
          this.loading = false
          this.$emit('upload')
          EventBus.$emit('file-upload')
          this.$bvModal.hide(this.id)
        })
        .catch(() => {
          this.loading = false
          // console.error(res)
        })
    },
    onClose () {
      this.files = []
      this.$emit('close')
    },
    show () {
      this.$refs['modal'].show()
    },
    hide () {
      this.$refs['modal'].hide()
    }
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

.border-blue-shadow {
  border: 1px solid $app-icon;
  box-shadow: 0 0.1em 0.25em 0.1em $app-icon;
}

.shadow-browse {
  border-radius: 5px;
}
</style>
