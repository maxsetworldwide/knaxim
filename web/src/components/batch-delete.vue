<script>
import FileService from '@/service/file'

export default {
  name: 'batch-delete',
  props: {
    files: {
      type: Array,
      required: true
    }
  },
  // TODO: The html for the delete list of files and its modal would go nicely
  // in a <template> with a slot; how do I pass inputEvents to the scoped slot.
  // ...the modal could be a b-modal component outside of the slot, instead of
  // a pragmatic msgBoxConfirm call.
  render () {
    return this.$scopedSlots.default({
      delete: this.delete,
      inputEvents: {
        click: (e) => { this.delete() }
      }
    })
  },

  methods: {
    async delete (files) {
      if (await this.$bvModal.msgBoxConfirm(this.createMsgBody(), {
        modalClass: 'modal-msg',
        title: 'The Following Files Will Be Deleted'
      })) {
        let error = []
        this.files.forEach(async file => {
          await FileService.erase({ fid: file.id }).catch(() => {
            error.push(file)
          })
        })

        if (!error.length) {
          this.$emit('delete-files')
        } else {
          // console.log(`Error: some files not deleted.`)
        }
      }
    },

    createMsgBody () {
      const h = this.$createElement
      return h('b-list-group', [
        h('b-list-group-item', 'Filename, Owner, Date'), ...this.files.map((file) => {
          return h('b-list-group-item',
            `${file.name}, ${file.own}, ${file.date.upload}`)
        })
      ])
    }
  }
}
</script>
