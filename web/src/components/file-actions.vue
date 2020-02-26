<template>
  <b-dropdown class="file-actions" no-caret variant="link" size="sm">
    <template v-slot:button-content>
      <svg class="more">
        <use href="../assets/app.svg#more" />
      </svg>
    </template>

    <b-dropdown-item href="#" :disabled="!fileSelected" @click="newFolder">
      <svg>
        <use href="../assets/app.svg#folder-2" />
      </svg>
      <span>Folder+</span>
    </b-dropdown-item>

    <b-dropdown-item
      href="#"
      :disabled="!fileSelected || activeFolders.length === 0"
      @click="removeFromFolder"
    >
      <b-icon-x-square />
      <span>Folder-</span>
    </b-dropdown-item>

    <b-dropdown-divider />

    <b-dropdown-item href="#" :disabled="!fileSelected" @click="share">
      <svg>
        <use href="../assets/app.svg#share" />
      </svg>
      <span>Share</span>
    </b-dropdown-item>

    <b-dropdown-item href="#" :disabled="!fileSelected" @click="adjustFavorite">
      <svg>
        <use href="../assets/app.svg#star" />
      </svg>
      <span v-if="isFavorite">UnFavorite</span>
      <span v-else>Favorite</span>
    </b-dropdown-item>

    <b-dropdown-item
      href="#"
      v-if="isTrash"
      :disabled="!fileSelected"
      @click="adjustFolder(false, '_trash_')"
    >
      <b-icon-arrow-counterclockwise />
      <span>Restore File</span>
    </b-dropdown-item>

    <b-dropdown-item v-if="singleFile" href="#" @click="downloadOriginal">
      <svg>
        <use href="../assets/app.svg#cloud" />
      </svg>
      <span>Download Original</span>
    </b-dropdown-item>

    <b-dropdown-item
      v-if="singleFile && !disableDownloadPDF"
      href="#"
      @click="downloadPdf"
    >
      <svg>
        <use href="../assets/app.svg#pdf" />
      </svg>
      <span>Download as PDF</span>
    </b-dropdown-item>

    <batch-delete
      v-if="!singleFile"
      :files="checkedFiles"
      :permanent="isTrash"
      #default="{ inputEvents }"
      v-on:delete-files="$emit('delete-files')"
    >
      <b-dropdown-item href="#" v-on="inputEvents" :disabled="!fileSelected">
        <svg>
          <use href="../assets/app.svg#bin" />
        </svg>
        <span>Delete</span>
      </b-dropdown-item>
    </batch-delete>

    <!-- <b-dropdown-divider/>
    <b-dropdown-item href="#" v-bind:disabled="fileSelected">
      <svg>
        <use href="../assets/app.svg#files"/>
      </svg>
      <span>Redact</span>
    </b-dropdown-item> -->
    <folder-modal
      ref="folderModal"
      id="newFolderModal"
      @new-folder="adjustFolder(true, $event)"
    />
    <share-modal
      ref="shareModal"
      id="file-list-share-modal"
      :files="checkedFiles"
    />
  </b-dropdown>
</template>

<script>
import BatchDelete from '@/components/batch-delete'
import FolderModal from '@/components/modals/folder-modal'
import ShareModal from '@/components/modals/share-modal'
import { PUT_FILE_FOLDER, REMOVE_FILE_FOLDER } from '@/store/actions.type'
import { mapGetters } from 'vuex'
import FileService from '@/service/file'

export default {
  name: 'file-actions',
  components: {
    BatchDelete,
    FolderModal,
    ShareModal
  },
  props: {
    checkedFiles: {
      type: Array,
      required: true
    },
    singleFile: {
      type: Boolean,
      default: false
    },
    disableDownloadPDF: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    newFolder () {
      this.showFolderModal()
    },
    adjustFavorite () {
      this.adjustFolder(!this.isFavorite, '_favorites_')
    },
    showShareModal () {
      this.$refs['shareModal'].show()
    },
    showFolderModal () {
      this.$refs['folderModal'].show()
    },
    adjustFolder (add, name) {
      this.checkedFiles.forEach(({ id: fid }) => {
        this.$store.dispatch(add ? PUT_FILE_FOLDER : REMOVE_FILE_FOLDER, {
          fid,
          name,
          group: this.activeGroup ? this.activeGroup.id : undefined
        })
      })
    },
    removeFromFolder () {
      const fileNames = this.checkedFiles.map(file => {
        return file.name
      })
      const folders = this.activeFolders
      const h = this.$createElement
      function msgBody () {
        return h('b-container', [
          h('b-row', [
            h('b-col', [
              h('h5', 'Files:'),
              h(
                'ul',
                fileNames.map(name => {
                  return h('li', name)
                })
              )
            ]),
            h('b-col', [
              h('h5', 'Folders:'),
              h(
                'ul',
                folders.map(folder => {
                  return h('li', folder)
                })
              )
            ])
          ])
        ])
      }
      this.$bvModal
        .msgBoxConfirm(msgBody(), {
          modalClass: 'modal-msg',
          title: 'Files will be removed from these folders:'
        })
        .then(val => {
          if (val) {
            folders.forEach(folder => {
              this.adjustFolder(false, folder)
            })
          }
        })
    },
    share () {
      this.showShareModal()
    },
    downloadOriginal () {
      if (this.checkedFiles[0]) {
        const fid = this.checkedFiles[0].id
        window.location.href = FileService.downloadURL({ fid })
      }
    },
    downloadPdf () {
      if (this.checkedFiles[0]) {
        const fid = this.checkedFiles[0].id
        window.location.href = FileService.viewURL({ fid })
      }
    }
  },
  computed: {
    fileSelected () {
      return this.checkedFiles.length > 0
    },
    trashFolder () {
      return this.getFolder('_trash_') || []
    },
    favoritesFolder () {
      return this.getFolder('_favorites_') || []
    },
    isFavorite () {
      return this.checkedFiles.reduce((acc, file) => {
        return acc && this.favoritesFolder.includes(file.id)
      }, true)
    },
    isTrash () {
      return this.checkedFiles.reduce((acc, file) => {
        return acc && this.trashFolder.includes(file.id)
      }, true)
    },
    ...mapGetters([
      'getFolder',
      'activeFolders',
      'activeGroup',
      'populateFiles'
    ])
  }
}
</script>

<style lang="scss">
.file-actions {
  .dropdown {
    height: 35px;
  }

  svg {
    width: 25px;
    height: 25px;
    margin-right: 15px;
  }

  .more {
    width: 40px;
    height: 40px;
  }
}
</style>
