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
<script>
import { PUT_FILE_FOLDER, LOAD_SERVER, DELETE_FILES } from '@/store/actions.type'
import { mapGetters } from 'vuex'

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

  computed: {
    ownedFiles () {
      return this.files.filter((f) => {
        return f.isOwned
      })
    },
    ...mapGetters({
      allOwned: 'ownedFiles'
    })
  },

  methods: {
    async delete () {
      if (await this.$bvModal.msgBoxConfirm(this.createMsgBody(), {
        modalClass: 'modal-msg',
        title: this.permanent ? 'The Following Files Will Be Deleted' : 'The Following Files Will Be Moved to Trash'
      })) {
        if (this.permanent) {
          let error = []
          try {
            await this.$store.dispatch(DELETE_FILES, { ids: this.ownedFiles.map(f => f.id) })
          } catch {
            // TODO: Handle Error
          }

          if (!error.length) {
            this.$emit('delete-files')
          }
        } else {
          let noerror = true
          try {
            await Promise.all(this.ownedFiles.map(async file => {
              await this.$store.dispatch(PUT_FILE_FOLDER, { fid: file.id, name: '_trash_', preventReload: true })
            })).finally(() => {
              this.$store.dispatch(LOAD_SERVER)
            })
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
        h('b-list-group-item', 'Filename, Upload Date'), ...this.ownedFiles.map((file) => {
          return h('b-list-group-item',
            `${file.name}, ${file.date.upload}`)
        })
      ])
    }
  }
}
</script>
