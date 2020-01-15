<template>
  <div class="h-100">
    <pdf-doc v-if="fileType === 'pdf'" :fileID="id" :acr="acr"/>
    <text-viewer v-else :fileName="name" :finalPage="sentenceCount" :acr="acr"/>
  </div>
</template>

<script>
import TextViewer from '@/components/text-viewer'
import PDFDoc from '@/components/pdf-doc'
import { TOUCH } from '@/store/mutations.type'
import FileService from '@/service/file'

export default {
  name: 'file-viewer',
  components: {
    'text-viewer': TextViewer,
    'pdf-doc': PDFDoc
  },
  props: {
    id: String,
    acr: String
  },
  data () {
    return {
      name: '',
      sentenceCount: 0
    }
  },
  computed: {
    fileType () {
      if (!this.name) {
        return 'txt'
      }
      return this.name.split('.').slice(-1)[0].toLowerCase()
    }
  },
  methods: {
    refresh () {
      // Add fileID to Recents
      this.$store.commit(TOUCH, this.id)
      FileService.info({ fid: this.id }).then(({ data }) => {
        const {
          count,
          file: { name }
        } = data
        this.name = name
        this.sentenceCount = count
      })
    }
  },
  created () {
    this.refresh()
  },
  watch: {
    id (to, from) {
      if (to !== from) {
        this.refresh()
      }
    }
  }
}
</script>

<style>

</style>
