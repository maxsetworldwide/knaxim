<!--
share-modal: window for sharing and unsharing files

props:
  id (req): id of the modal
  files (req): array of file OBJECTS to operate upon

events:
  'close': emitted when modal is closed
-->
<template>
  <b-modal
    :id="id"
    ref="modal"
    @hidden="onClose"
    centered
    title="Files to Share"
    hide-footer
    size="lg"
    content-class="modal-style"
  >
    <b-container>
      <!-- <b-row align-h="center" v-if="!hideList">
        <h4>Files to share:</h4>
      </b-row> -->
      <b-row align-h="center" v-if="!hideList">
        <ul class="file-list">
          <li v-for="file in files" :key="file.id">{{ file.name }}</li>
        </ul>
      </b-row>
      <b-form @submit.prevent="share">
        <b-row align-h="center">
          <b-col cols="6">
            <b-form-input
              :state="validName"
              v-model="inputName"
              placeholder="Enter username or team name"
              autofocus
              debounce="400"
            />
            <b-form-invalid-feedback>Name not found!</b-form-invalid-feedback>
          </b-col>
          <b-col class="text-center" cols="2">
            <b-spinner v-if="shareLoading" />
            <b-button v-else :disabled="!validName" type="submit"
              >Share</b-button
            >
          </b-col>
        </b-row>
      </b-form>
      <b-row class="mt-3" align-h="center">
        <h4>Current Viewers:</h4>
      </b-row>
      <b-row class="w-75 mx-auto viewer-list" align-h="center">
        <div v-if="objIsEmpty(viewers)">
          <span v-if="!loadingViewers"
            >There are currently no viewers for these files.</span
          >
          <div v-else>
            <b-spinner class="align-middle"></b-spinner>
            <strong>Loading...</strong>
          </div>
        </div>
        <!-- Relying on b-col's behavior of wrapping when #cols exceeds 12 -->
        <b-col cols="3" v-for="(name, uid) in viewers" :key="uid">
          <share-viewer
            :uid="uid"
            :name="name"
            :fids="fileIDs"
            @stop-share="getViewers"
          />
        </b-col>
      </b-row>
    </b-container>
  </b-modal>
</template>

<script>
import PermissionService from '@/service/permission'
import UserService from '@/service/user'
import GroupService from '@/service/group'
import ShareViewer from '@/components/share-viewer'
import Vue from 'vue'

export default {
  name: 'share-modal',
  components: {
    ShareViewer
  },
  props: {
    id: {
      type: String,
      required: true
    },
    files: {
      type: Array,
      required: true
    },
    hideList: Boolean
  },
  data () {
    return {
      shareLoading: false,
      inputName: '',
      nameID: '',
      validName: null,
      viewers: {},
      loadingViewers: false
    }
  },
  created () {
    this.getViewers()
  },
  watch: {
    files () {
      this.getViewers()
    },
    inputName () {
      this.validateName()
    }
  },
  computed: {
    fileIDs () {
      return this.files.map(file => {
        return file.id
      })
    }
  },
  methods: {
    objIsEmpty (obj) {
      for (let key in obj) {
        if (obj.hasOwnProperty(key)) {
          return false
        }
      }
      return true
    },
    validateName () {
      const currName = this.inputName
      if (currName.length === 0) {
        this.validName = null
        return
      }
      UserService.lookup({ name: currName })
        .then(({ data }) => {
          if (!data.id) {
            throw new Error('name is not user')
          }
          this.validName = !!data.id && this.inputName === currName
          if (this.validName) {
            this.nameID = data.id
          } else {
            this.nameID = ''
          }
        })
        .catch(() => {
          GroupService.lookup({ name: currName })
            .then(({ data }) => {
              this.validName = !!data.id && this.inputName === currName
              if (this.validName) {
                this.nameID = data.id
              } else {
                this.nameID = ''
              }
            })
            .catch(() => {
              this.validName = false
            })
        })
    },
    getViewers () {
      if (this.files.length === 0) {
        this.viewers = {}
        return
      }
      this.loadingViewers = true
      let intersectNames = []
      Promise.all(
        this.fileIDs.map(id => {
          return PermissionService.permissions({ id })
        })
      )
        .then(resArray => {
          const viewerLists = resArray.map(({ data }) => {
            if (data.permission.view) {
              return data.permission.view
            } else {
              return []
            }
          })
          intersectNames = viewerLists.reduce((acc, curr) => {
            return curr.filter(val => {
              return acc.indexOf(val) > -1
            })
          })
          this.viewers = {}
          if (intersectNames.length === 0) {
            this.loadingViewers = false
            return
          }
          let counter = Promise.resolve(intersectNames.length)
          intersectNames.forEach(id => {
            UserService.info({ id }).then(({ data }) => {
              if (data.name) {
                Vue.set(this.viewers, id, data.name)
              }
              counter.then(count => {
                count -= 1
                if (count === 0) {
                  this.loadingViewers = false
                }
                return count
              })
            })
            GroupService.info({ gid: id }).then(({ data }) => {
              if (data.name) {
                Vue.set(this.viewers, id, data.name)
              }
              counter.then(count => {
                count -= 1
                if (count === 0) {
                  this.loadingViewers = false
                }
                return count
              })
            })
          })
        })
        .finally(() => {
          this.loadingViewers = false
        })
    },
    share () {
      if (this.validName) {
        this.shareLoading = true
        Promise.all(
          this.fileIDs.map(fid => {
            return PermissionService.share({ id: fid, targets: this.nameID })
          })
        )
          .then(res => {
            this.getViewers()
          })
          .catch(res => {
            // console.log('share error:', res)
          })
          .finally(() => {
            this.shareLoading = false
          })
      }
    },
    onClose () {
      this.$emit('close')
    },
    show () {
      this.$refs['modal'].show()
    }
  }
}
</script>

<style scoped lang="scss">
.file-list {
  height: 90px;
  width: 60%;
  overflow: auto;
  border: 2px solid $app-clr2;
  border-radius: 8px;
}

.viewer-list {
  height: 100px;
  overflow: auto;
  border: 2px solid $app-clr2;
  border-radius: 8px;
}

button {
  @extend %pill-buttons;
  border: 1px solid $app-clr;
  width: 100%;
}

::v-deep .modal-style {
  @extend %modal-corners;
}
</style>
