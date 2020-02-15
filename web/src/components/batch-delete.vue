<script>
import FileService from '@/service/file'
import { PUT_FILE_FOLDER } from '@/store/actions.type'

export default {
  name: 'batch-delete',
  props: {
    files: {
      type: Array,
      required: true
    },
    permanent: Boolean
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
        title: this.permanent ? 'The Following Files Will Be Deleted' : 'The Following Files Will Be Moved to Trash'
      })) {
        if (this.permanent) {
          let error = []
          try {
            await Promise.all(this.files.map(async file => {
              await FileService.erase({ fid: file.id }).catch(() => {
                error.push(file)
              })
            }))
          } catch {}

          if (!error.length) {
            this.$emit('delete-files')
          }
        } else {
          let noerror = true
          try {
            await Promise.all(this.files.map(async file => {
              await this.$store.dispatch(PUT_FILE_FOLDER, { fid: file.id, name: '_trash_' })
            }))
          } catch {
            noerror = false
          }
          if (noerror) {
            this.$emit('delete-files')
          }
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
