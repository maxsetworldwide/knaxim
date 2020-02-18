<template>
  <!-- prompt Login -->

  <!-- No files exist -->
  <div v-if="promptUpload && !loading" class="empty">
    <h1>No files!</h1>
    <b-button v-b-modal="'file-list-upload'">
      <h3>Add file?</h3>
    </b-button>
    <upload-modal id="file-list-upload"/>
  </div>

  <div v-else>
    <p v-if="folderFilters.length > 0">
      Open Folders: <span v-for="fold in folderFilters" :key="fold">{{ fold }} <span @click="closeFolder(fold)" class="removeFolder">X</span> </span>
    </p>
     <file-table
        :files="fileids"
        :folders="folderRows"
        :busy="loading"
        @selection="onCheck"
        @open-folder="openFolder"
        @open="open"
      >
       <template #action>
         <file-list-batch :checkedFiles="selected"
           :removeFavorite="src === 'favorites'"
           :fileSelected="checked.length > 0"
           @favorite="adjustFavorite"
           @add-folder="showFolderModal"
           @share-file="showShareModal"
         />
         <folder-modal
           ref="folderModal"
           id="newFolderModal"
           @new-folder="createFolder"
         />
         <share-modal
           ref="shareModal"
           id="file-list-share-modal"
           :files="selected"
         />
       </template>
     </file-table>
  </div>

</template>

<script>
import UploadModal from '@/components/modals/upload-modal'
import FolderModal from '@/components/modals/folder-modal'
import ShareModal from '@/components/modals/share-modal'
import FileListBatch from '@/components/file-list-batch'
import FileTable from '@/components/file-table'
import { LOAD_FOLDERS, PUT_FILE_FOLDER, REMOVE_FILE_FOLDER, GET_USER } from '@/store/actions.type'
import { ACTIVATE_FOLDER, DEACTIVATE_FOLDER } from '@/store/mutations.type'
import { mapGetters } from 'vuex'

export default {
  name: 'file-list',
  components: {
    UploadModal,
    FileListBatch,
    FolderModal,
    ShareModal,
    FileTable
  },
  props: {
    src: String
  },
  data () {
    return {
      checked: [], // Only manipulated via onCheck method
      folderFilters: []
    }
  },
  created () {
    // EventBus.$on(['file-upload', 'url-upload'], this.refresh)
    // this.refresh()
  },
  computed: {
    files () {
      if (this.src === 'recents') {
        return this.recentFiles || []
      } else if (this.src === 'favorites') {
        return this.folders['_favorites_'] || []
      } else if (this.src === 'shared') {
        return this.sharedFiles
      } else if (this.src === 'owned') {
        return this.ownedFiles
      }
      // console.log(this.ownedFiles)
      // console.log(this.sharedFiles)
      return [...this.ownedFiles, ...this.sharedFiles]
    },
    fileids () {
      return this.files.filter(id => {
        return this.activeFolders.reduce((acc, name) => {
          return acc && this.folders[name].indexOf(id) > -1
        }, true)
      })
    },
    folderRows () {
      let rows = []
      for (let name in this.folders) {
        if (name !== '_favorites_' && this.activeFolders.reduce((acc, active) => {
          return acc && active !== name
        }, true)) {
          rows.push(name)
        }
      }
      return rows
    },
    promptUpload () {
      if (this.src) {
        return false
      }
      return this.files.length === 0
    },
    selected () {
      return this.populateFiles(this.checked)
    },
    ...mapGetters(['isAuthenticated', 'ownedFiles', 'sharedFiles', 'recentFiles', 'folders', 'activeFolders', 'activeGroup', 'loading', 'populateFiles'])
  },
  methods: {
    onCheck (rows) {
      this.checked = rows
    },
    open (id) {
      this.$router.push(`/file/${id}`)
    },
    openFolder (name) {
      this.$store.commit(ACTIVATE_FOLDER, name)
    },
    closeFolder (name) {
      this.$store.commit(DEACTIVATE_FOLDER, name)
    },
    // file actions
    adjustFavorite (add) {
      this.checked.forEach((fid) => {
        this.$store.dispatch(add ? PUT_FILE_FOLDER : REMOVE_FILE_FOLDER, {
          fid,
          name: '_favorites_',
          group: (this.activeGroup.id || undefined)
        })
      })
    },
    showShareModal () {
      this.$refs['shareModal'].show()
    },
    showFolderModal () {
      this.$refs['folderModal'].show()
    },
    createFolder (name) {
      this.checked.forEach((fid) => {
        this.$store.dispatch(
          PUT_FILE_FOLDER,
          { fid, name, group: (this.activeGroup.id || undefined) }
        )
      })
    }
  },
  watch: {
    gid (n, o) {
      if (n !== o) {
        this.refresh()
      }
    }
  },
  mounted () {
    this.$store.dispatch(GET_USER, {})
    this.$store.dispatch(LOAD_FOLDERS, {})
  }
}
</script>

<style lang="scss" scoped>

  .divide {
    margin-top: 2px;
    margin-bottom: 2px;
    height: 1px;
    width: 100%;
    border-top: 1px solid gray;
  }

  .head-row {
    margin-top: 10px;
  }

  .empty {
    text-align: center;
    margin-top: 10%;
    button {
      background-color: white;
      border-radius: 10px;
      border: 0px;
      width: 160px;
      height: 80px;
      color: rgb(46, 46, 46);
    }

    button:hover {
      background-color: rgb(150, 182, 252);
      color: rgb(46, 46, 46);
    }
  }

  .file-name {
    cursor: pointer;
    color: $app-clr1;

    &:hover {
      text-decoration: underline;
    }
  }
  svg {
    width: 100%;
    height: 40px;
  }

  .action-column {
    width: 5%;
  }

  .name-column {
    min-width: 30%;
  }

  .expand-column {
    width: 8%;
  }

  .triangle {
    fill: gray;

    &:hover {
      fill: $app-bg4;
    }
  }

  .removeFolder {
    color: red;
  }

  .removeFolder:hover {
    text-decoration: underline;
    cursor: pointer;
  }
</style>
