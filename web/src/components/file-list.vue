<template>
  <!-- prompt Login -->

  <!-- No files exist -->
  <div v-if="promptUpload && !loading" class="empty">
    <h1>No files!</h1>
    <b-button v-b-modal="'file-list-upload'">
      <h3>Add file?</h3>
    </b-button>
    <upload-modal id="file-list-upload" />
  </div>

  <div v-else>
    <!--
       - <p v-if="activeFolders.length > 0">
       -   Open Folders:
       -   <span v-for="fold in activeFolders" :key="fold"
       -     >{{ fold }}
       -     <span @click="closeFolder(fold)" class="removeFolder">X</span>
       -   </span>
       - </p>
       -->
    <file-table
      :files="files"
      :busy="loading"
      :trashmode="src === 'trash'"
      @selection="onCheck"
      @open="open"
    >
      <template #action>
        <file-actions :checkedFiles="selected" />
      </template>
    </file-table>
  </div>
</template>

<script>
import UploadModal from '@/components/modals/upload-modal'
import FileActions from '@/components/file-actions'
import FileTable from '@/components/file-table'
import { LOAD_FOLDERS, GET_USER } from '@/store/actions.type'
import { SET_ACTIVE_FILES } from '@/store/mutations.type'
import { mapGetters, mapMutations } from 'vuex'

export default {
  name: 'file-list',
  components: {
    UploadModal,
    FileActions,
    FileTable
  },
  props: {
    src: String
  },
  data () {
    return {
      checked: [] // Only manipulated via onCheck method
    }
  },
  created () {
    // EventBus.$on(['file-upload', 'url-upload'], this.refresh)
    // this.refresh()
  },
  computed: {
    files () {
      let trashFolder = this.folders['_trash_'] || []
      function filterTrash (ids) {
        return ids.filter(id => {
          return trashFolder.reduce((a, i) => {
            return a && i !== id
          }, true)
        })
      }
      if (this.src === 'recents') {
        return filterTrash(this.recentFiles || [])
      } else if (!this.activeGroup && this.src === 'favorites') {
        return filterTrash(this.folders['_favorites_'] || [])
      } else if (!this.activeGroup && this.src === 'shared') {
        return filterTrash(this.sharedFiles)
      } else if (!this.activeGroup && this.src === 'owned') {
        return filterTrash(this.ownedFiles)
      } else if (!this.activeGroup && this.src === 'trash') {
        return trashFolder
      }
      // console.log(this.ownedFiles)
      // console.log(this.sharedFiles)
      return filterTrash([...this.ownedFiles, ...this.sharedFiles])
    },
    /*
     * fileids () {
     *   return this.files.filter(id => {
     *     return this.activeFolders.reduce((acc, name) => {
     *       return acc && this.folders[name].indexOf(id) > -1
     *     }, true)
     *   })
     * },
     */
    /*
     * folderRows () {
     *   let rows = []
     *   for (let name in this.folders) {
     *     if (
     *       name !== '_favorites_' &&
     *       name !== '_trash_' &&
     *       this.activeFolders.reduce((acc, active) => {
     *         return acc && active !== name
     *       }, true)
     *     ) {
     *       rows.push(name)
     *     }
     *   }
     *   return rows
     * },
     */
    promptUpload () {
      if (this.src || this.activeGroup || this.loading) {
        return false
      }
      return this.files.length === 0
    },
    selected () {
      return this.populateFiles(this.checked)
    },
    ...mapGetters([
      'isAuthenticated',
      'ownedFiles',
      'sharedFiles',
      'recentFiles',
      'folders',
      // 'activeFolders',
      'activeGroup',
      'loading',
      'populateFiles'
    ])
  },
  methods: {
    onCheck (rows) {
      this.checked = rows
    },
    open (id) {
      this.$router.push(`/file/${id}`)
    },
    ...mapMutations({
      currentFiles: SET_ACTIVE_FILES
    })
  },
  watch: {
    gid (n, o) {
      if (n !== o) {
        this.refresh()
      }
    },
    files (n) {
      this.currentFiles(n)
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
/*
 *
 * .removeFolder {
 *   color: red;
 * }
 *
 * .removeFolder:hover {
 *   text-decoration: underline;
 *   cursor: pointer;
 * }
 */
</style>
