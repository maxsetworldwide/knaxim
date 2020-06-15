<template>
  <b-container class="image-viewer">
    <b-row class="image-toolbar">
      <b-col class="d-none d-md-flex" md="1">
        <b-button v-if="navigation.prev" @click="openPrev" class="min-width-1em">
          <b-icon icon="arrow-bar-left" class="icon" />
        </b-button>
      </b-col>
      <b-col offset="2" cols="6" offset-md="3" md="4">
        <h4 class="title text-center">{{ fileInfo.name }}</h4>
      </b-col>
      <b-col cols="2" md="1">
        <file-actions
          fileSelected
          singleFile
          disableDownloadPDF
          :checkedFiles="[fileInfo]"
        />
      </b-col>
      <b-col class="d-none d-md-flex" offset-md="2" md="1">
        <b-button v-if="navigation.next" @click="openNext" class="min-width-1em">
          <b-icon icon="arrow-bar-right" class="icon" />
        </b-button>
      </b-col>
    </b-row>
    <b-row class="content-row h-100" align-h="center">
      <img
        class="image-content"
        :src="srcURL"
        :alt="fileInfo.name"
        @error="$emit('no-image')"
      />
    </b-row>
  </b-container>
</template>

<script>
import FileActions from '@/components/file-actions'
import FileService from '@/service/file'
import { mapGetters } from 'vuex'

export default {
  name: 'image-viewer',
  components: {
    FileActions
  },
  props: {
    id: {
      type: String,
      required: true
    }
  },
  computed: {
    srcURL () {
      return FileService.downloadURL({ fid: this.id })
    },
    fileInfo () {
      return this.populateFiles(this.id)
    },
    navigation () {
      let nav = {
        prev: null,
        next: null
      }
      let indx = this.activeFiles.findIndex(af => af === this.id)
      if (indx > 0) {
        nav.prev = this.activeFiles[indx - 1]
      }
      if (indx >= 0 && indx < this.activeFiles.length - 1) {
        nav.next = this.activeFiles[indx + 1]
      }
      return nav
    },
    ...mapGetters(['populateFiles', 'activeFiles'])
  },
  methods: {
    openPrev () {
      if (this.navigation.prev) {
        this.$router.push(`/file/${this.navigation.prev}`)
      }
    },
    openNext () {
      if (this.navigation.next) {
        this.$router.push(`/file/${this.navigation.next}`)
      }
    }
  }
}
</script>

<style scoped lang="scss">
.title {
  width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.image-viewer {
  height: 90%;
}

.image-toolbar {
  height: 10%;
}

.content-row {
  overflow: auto;
}

.image-content {
  object-fit: contain;
}

button {
  background-color: white;
  border-radius: 10px;
  border: 0px;
  width: 100%;
  height: 30px;
  color: rgb(46, 46, 46);
}

button:hover {
  background-color: rgb(150, 182, 252);
  color: rgb(46, 46, 46);
}
</style>
