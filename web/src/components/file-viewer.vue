<template>
  <div class="h-100">
    <pdf-doc
      v-if="viewExists"
      :fileID="id"
      :acr="acr"
      @no-view="viewExists = false"
    />
    <image-viewer v-else-if="imageTypes.includes(fileType)" :id="id" />
    <text-viewer
      v-else
      :fileName="name"
      :finalPage="sentenceCount"
      :acr="acr"
    />
  </div>
</template>

<script>
import TextViewer from '@/components/text-viewer'
import PDFDoc from '@/components/pdf/pdf-doc'
import ImageViewer from '@/components/image-viewer'
import { TOUCH } from '@/store/mutations.type'
import { GET_FILE } from '@/store/actions.type'

export default {
  name: 'file-viewer',
  components: {
    'text-viewer': TextViewer,
    'pdf-doc': PDFDoc,
    'image-viewer': ImageViewer
  },
  props: {
    id: String,
    acr: String
  },
  data () {
    return {
      name: '',
      sentenceCount: 0,
      viewExists: true,
      imageTypes: [
        'jpg',
        'jpeg',
        'jfif',
        'pjpeg',
        'pjp',
        'png',
        'gif',
        'apng',
        'bmp',
        'ico',
        'cur',
        'svg',
        'webp'
      ]
    }
  },
  computed: {
    fileType () {
      if (!this.name) {
        return 'txt'
      }
      return this.name
        .split('.')
        .slice(-1)[0]
        .toLowerCase()
    }
  },
  methods: {
    refresh () {
      // Add fileID to Recents
      this.$store.commit(TOUCH, this.id)
      this.$store.dispatch(GET_FILE, this).then(data => {
        const { count, name } = data
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

<style></style>
