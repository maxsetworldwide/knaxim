<template>
  <!-- prompt Login -->
  <div v-if="!isAuthenticated" class="empty">
    <h1>You aren't logged in!</h1>
    <b-button @click="showAuth">
      <h3>Login</h3>
    </b-button>
  </div>

  <!-- No files exist -->
  <div v-else-if="promptUpload" class="empty">
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
    <b-table
      ref="fileTable"
      striped
      hover
      selectable
      :items="rows"
      :fields="columnHeaders"
      :busy="loading"
      :sort-compare="sortCompare"
      @row-selected="onCheck"
    >
      <template v-slot:table-colgroup="scope">
        <col
          v-for="field in scope.fields"
          :key="field.key"
          :class="field.class"
        >
      </template>
      <template v-slot:head(expand)="col">
        <svg @click.stop="expandAll">
          <use href="../assets/app.svg#expand-tri" class="triangle"/>
        </svg>
      </template>
      <template v-slot:table-busy>
        <div class="text-center">
          <b-spinner class="align-middle"></b-spinner>
          <strong>Loading...</strong>
        </div>
      </template>
      <template v-slot:head(select)>
        <b-checkbox v-model="selectAllMode"/>
      </template>
      <template v-slot:cell(select)="{ rowSelected }">
        <template v-if="rowSelected">
          <span aria-hidden="true">&check;</span>
        </template>
        <template v-else>
          <span aria-hidden="true">&nbsp;</span>
        </template>
      </template>
      <template v-slot:head(action)>
        <file-list-batch :checkedFiles="selected"
          :removeFavorite="src === 'favorites'"
          :fileSelected="selectAllMode"
          @favorite="adjustFavorite"
          @add-folder="showFolderModal"
          @delete-files="refresh"
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
      <template v-slot:cell(name)="data">
        <span v-if="data.item.isFolder" class="file-name" @click.prevent.stop="openFolder(data.value)">{{ data.value }}</span>
        <span v-else class="file-name" @click="open(data.item.id)">{{ data.value }}</span>
      </template>
      <template v-slot:cell(expand)="row">
        <svg v-if="row.item.preview" @click.stop="row.toggleDetails">
          <use href="../assets/app.svg#expand-tri" class="triangle"/>
        </svg>
      </template>
      <template v-slot:row-details="row">
        <span>{{ row.item.preview }}</span>
      </template>
      <template v-slot:cell(action)="data">
        <file-icon :extention="(data.item.ext || '')" :folder="data.item.isFolder" :webpage="!!data.item.url"/>
      </template>
    </b-table>
  </div>

</template>

<script>
import uploadModal from '@/components/modals/upload-modal'
import FolderModal from '@/components/modals/folder-modal'
import ShareModal from '@/components/modals/share-modal'
import fileListBatch from '@/components/file-list-batch'
import fileIcon from '@/components/file-icon'
import { LOAD_FOLDERS, PUT_FILE_FOLDER, REMOVE_FILE_FOLDER, GET_USER } from '@/store/actions.type'
import { mapGetters } from 'vuex'
import FileService from '@/service/file'
import UserService from '@/service/user'
import GroupService from '@/service/group'
import FolderService from '@/service/folder'
import { EventBus } from '@/main'
import Vue from 'vue'
import { humanReadableSize, humanReadableTime } from '@/plugins/utils'

export default {
  name: 'file-list',
  components: {
    uploadModal,
    fileListBatch,
    FolderModal,
    ShareModal,
    fileIcon
  },
  props: {
    src: String
  },
  data () {
    return {
      checked: [], // Only manipulated via onCheck method
      fileSet: {
        owned: {},
        shared: {}
      },
      loading: true,
      recents: [],
      ownerNames: {},
      folderFilters: [],
      columnHeaders: [
        {
          key: 'select'
        },
        {
          key: 'action',
          class: 'action-column'
        },
        {
          key: 'name',
          class: 'name-column',
          sortable: true
        },
        {
          key: 'expand',
          label: '',
          class: 'expand-column'
        },
        {
          key: 'owner',
          sortable: true
        },
        {
          key: 'date',
          sortable: true
        },
        {
          key: 'size',
          sortable: true
        }
      ]
    }
  },
  created () {
    EventBus.$on(['file-upload', 'url-upload'], this.refresh)
    this.refresh()
  },
  computed: {
    gid () {
      if (!this.activeGroup) {
        return ''
      }
      return this.activeGroup.id
    },
    files () {
      if (this.src === 'recents') {
        return this.recents
      } else if (this.src === 'favorites') {
        return (this.folders['_favorites_'] || []).map(
          function (id) {
            return this[id]
          },
          { ...this.fileSet.owned, ...this.fileSet.shared }
        )
      } else if (this.src === 'shared') {
        return Object.values(this.fileSet.shared)
      } else if (this.src === 'owned') {
        return Object.values(this.fileSet.owned)
      }
      return [...Object.values(this.fileSet.owned), ...Object.values(this.fileSet.shared)]
    },
    filterFiles () {
      return this.files.filter(file => {
        return this.folderFilters.map(name => this.folders[name]).reduce((acc, content) => {
          return acc && content.reduce((accu, id) => {
            return accu || id === file.id
          }, false)
        }, true)
      })
    },
    fileRows () {
      return this.filterFiles.map((file) => {
        let splitname = [file.name, '']
        if (file.name.split && !file.url) {
          let splits = file.name.split('.')
          if (splits.length > 1) {
            splitname[0] = splits.slice(0, -1).join('.')
            splitname[1] = splits[splits.length - 1]
          }
        }
        return {
          id: file.id,
          url: file.url,
          isFolder: false,
          name: splitname[0],
          ext: splitname[1],
          owner: this.ownerNames[file.own],
          size: file.size && humanReadableSize(file.size),
          sizeInt: file.size,
          date: file.date && humanReadableTime(file.date.upload),
          dateInt: Date.parse(file.date.upload),
          preview: file.preview,
          _showDetails: (file._showDetails || false)
        }
      })
    },
    folders () {
      if (!this.gid) {
        return this.$store.state.folder.user
      } else {
        return this.$store.state.folder.group[this.gid]
      }
    },
    folderRows () {
      var rows = []
      var id = 0
      for (const name in this.folders) {
        if (name !== '_favorites_' && this.folderFilters.reduce((acc, f) => {
          return acc && f !== name
        }, true)) {
          rows.push({
            name,
            isFolder: true,
            id
          })
          id++
        }
      }
      return rows
    },
    rows () {
      return [ ...this.folderRows, ...this.fileRows ]
    },
    promptUpload () {
      if (this.src) {
        return false
      }
      return this.files.length === 0
    },
    selected () {
      let allFiles = { ...this.fileSet.owned, ...this.fileSet.shared }
      this.recents.forEach((file) => {
        if (!allFiles[file.id]) {
          allFiles[file.id] = file
        }
      })
      return this.checked.map((fid) => allFiles[fid])
    },
    selectAllMode: {
      get () {
        return this.checked.length > 0
      },
      set (newValue) {
        if (!newValue && this.checked.length > 0) {
          this.unselectAll()
        }
        if (newValue && this.checked.length === 0) {
          this.selectAll()
        }
      }
    },
    anyRowExpanded () {
      return this.fileRows.reduce((acc, row) => {
        return acc || row._showDetails
      }, false)
    },
    ...mapGetters(['isAuthenticated', 'fileMap', 'activeGroup'])
  },
  methods: {
    refresh () {
      this.loading = true
      Promise.all([
        // Recent Files (:possibly from different groups)
        ...this.$store.state.recents.files.map((fid, indx) => {
          this.recents.push({
            id: 'loading',
            name: 'Loading...',
            count: 0,
            size: 0
          })
          return FileService.info({ fid }).then((res) => {
            this.populateOwnerName(res.data.file.own)
            Vue.set(this.recents, indx, {
              ...res.data.file,
              count: res.data.count,
              size: res.data.size
            })
          })
        }),

        // Owned File List
        FileService.list({ gid: (this.gid || undefined) }).then(
          (res) => {
            let fSet = {}
            for (let id in res.data.files) {
              let finfo = res.data.files[id]
              this.populateOwnerName(finfo.file.own)
              Vue.set(fSet, id, {
                ...finfo.file,
                count: finfo.count,
                size: finfo.size
              })
            }
            Vue.set(this.fileSet, 'owned', fSet)
          }
        ),

        // Shared File List
        FileService.list({ gid: (this.gid || undefined), shared: true }).then(
          (res) => {
            let fSet = {}
            for (let id in res.data.files) {
              let finfo = res.data.files[id]
              this.populateOwnerName(finfo.file.own)
              Vue.set(fSet, id, {
                ...finfo.file,
                count: finfo.count,
                size: finfo.size
              })
            }
            Vue.set(this.fileSet, 'shared', fSet)
          }
        ),
        FolderService.list({ group: (this.gid || undefined) }).then(
          ({ data }) => {
            data.folders = (data.folders || [])
            return Promise.all(data.folders.map((name) => {
              return FolderService.info({ name, group: (this.gid || undefined) }).then(({ data }) => {
                Vue.set(this.folders, name, data.files)
              })
            }))
          }
        )
      ]).then(() => {
        var list = this.recents.map(({ id }, indx) => {
          return FileService.slice({ fid: id, start: 0, end: 3 }).then(({ data }) => {
            const preview = data.lines.reduce((acc, line) => {
              return acc + line.Content[0]
            }, '')
            Vue.set(this.recents[indx], 'preview', preview)
          })
        })
        for (const id in this.fileSet.owned) {
          list.push(
            FileService.slice({ fid: id, start: 0, end: 3 })
              .then(({ data: { lines = [] } }) => {
                const preview = lines.reduce((acc, line) => {
                  return acc + line.Content[0]
                }, '')
                Vue.set(this.fileSet.owned[id], 'preview', preview)
              })
          )
        }
        for (const id in this.fileSet.shared) {
          list.push(
            FileService.slice({ fid: id, start: 0, end: 3 }).then(({ data }) => {
              const preview = data.lines.reduce((acc, line) => {
                return acc + line.Content[0]
              }, '')
              Vue.set(this.fileSet.shared[id], 'preview', preview)
            })
          )
        }
        return Promise.all(list)
      }).then(() => {
        this.loading = false
      }).catch((res) => {
        this.loading = false
      })
    },
    showAuth () {
      this.$router.push('/login')
    },
    sortCompare (a, b, key) {
      if (key === 'date') {
        return a.dateInt - b.dateInt
      } else if (key === 'size') {
        return a.sizeInt - b.sizeInt
      } else {
        return null
      }
    },
    onCheck (rows) {
      this.checked = rows.filter(row => !row.isFolder).map(row => row.id)
    },
    selectAll () {
      this.$refs.fileTable.selectAllRows()
    },
    unselectAll () {
      this.$refs.fileTable.clearSelected()
    },
    expandAll () {
      const expand = !this.anyRowExpanded
      for (const id in this.fileSet.owned) {
        Vue.set(this.fileSet.owned[id], '_showDetails', expand)
      }
      for (const id in this.fileSet.shared) {
        Vue.set(this.fileSet.shared[id], '_showDetails', expand)
      }
      this.recents.forEach((file, i, arr) => {
        Vue.set(arr[i], '_showDetails', expand)
      })
    },
    async populateOwnerName (id, overwrite) {
      if (overwrite || !this.ownerNames[id]) {
        Vue.set(this.ownerNames, id, 'loading...')
        let userprom = UserService.info({ id })
          .then((res) => {
            return res.data.name
          }, () => {
            return ''
          })
        let groupname = await GroupService.info({ gid: id })
          .then((res) => {
            return res.data.name
          }, () => {
            return ''
          })
        let username = await userprom
        if (username) {
          Vue.set(this.ownerNames, id, username)
        } else {
          Vue.set(this.ownerNames, id, groupname)
        }
      }
    },
    open (id) {
      this.$router.push(`/file/${id}`)
    },
    openFolder (name) {
      this.folderFilters.push(name)
    },
    closeFolder (name) {
      this.folderFilters = this.folderFilters.filter(f => f !== name)
    },
    // file actions
    adjustFavorite (add) {
      Promise.all(this.checked.map((fid) => {
        return this.$store.dispatch(add ? PUT_FILE_FOLDER : REMOVE_FILE_FOLDER, { fid, name: '_favorites_', group: (this.gid || undefined) })
      }))
        .then(() => { this.loading = false })
      this.loading = true
    },
    showShareModal () {
      this.$refs['shareModal'].show()
    },
    showFolderModal () {
      this.$refs['folderModal'].show()
    },
    createFolder (name) {
      Promise.all(this.checked.map((fid) => {
        return this.$store.dispatch(
          PUT_FILE_FOLDER,
          { fid, name, group: (this.gid || undefined) }
        )
      }))
        .then(() => {
          this.refresh()
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
