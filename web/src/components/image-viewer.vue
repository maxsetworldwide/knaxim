<template>
  <b-container class="h-100">
    <b-row>
      <b-col offset="3" cols="6" offset-md="4" md="4">
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
    </b-row>
    <b-row class="content-row h-100" align-h="center">
      <img class="image-content" :src="srcURL" :alt="fileInfo.name" />
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
    ...mapGetters(['populateFiles'])
  }
}
</script>

<style>
.title {
  width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.content-row {
  overflow: auto;
}

.image-content {
  object-fit: contain;
}
</style>
