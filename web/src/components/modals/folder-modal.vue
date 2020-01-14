<template>
  <b-modal
    ref="folderModal"
    :id="id"
    centered
    hide-footer
    hide-header
    @hidden="onClose"
    content-class="modal-style"
    >
    <b-form @submit.prevent="create">
      <input type="text" v-model="name" placeholder="Enter Name..."/>
      <b-button class="shadow-sm" type="submit">Create</b-button>
    </b-form>
  </b-modal>
</template>
<script>

export default {
  name: 'folder-modal',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      name: ''
    }
  },
  methods: {
    create () {
      this.$emit('new-folder', this.name)
      this.hide()
    },
    onClose () {
      this.$emit('close')
    },
    show () {
      this.$refs['folderModal'].show()
    },
    hide () {
      this.$refs['folderModal'].hide()
    }
  }
}
</script>
