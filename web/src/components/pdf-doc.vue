<template>
  <div>
    <b-container>
      <b-row align-h="center">
        <b-col cols="1">
          <b-button @click="increaseScale">
            <svg>
              <use href="@/assets/app.svg#zoom-in"></use>
            </svg>
          </b-button>
        </b-col>
        <b-col cols="1">
          <b-button @click="decreaseScale">
            <svg>
              <use href="@/assets/app.svg#zoom-out"></use>
            </svg>
          </b-button>
        </b-col>
        <b-col cols="2">
          <input
            value="1"
            @input="pageInput"
            min="1"
            :max="pages.length"
            type="number"
          />
          <span> / {{ pages.length }}</span>
        </b-col>

        <b-col>
          <file-list-batch fileSelected="Uh-huh" :checkedFiles="mockFileArray"
            @favorite="adjustFavorite"
            @add-folder="showFolderModal"
            @share-file="showShareModal"
          />
        </b-col>
      </b-row>
    </b-container>
    <div class="pdf-doc" :ref="'pdfdoc'">
      <PDFPage
        v-for="page in pages"
        v-bind="{ page, scale, scrollTop, clientHeight }"
        :key="page.pageNumber"
        ref="pages"
      />
    </div>

    <folder-modal
      ref="folderModal"
      id="newFolderModal"
      @new-folder="createFolder"
    />
    <share-modal
      ref="shareModal"
      id="file-list-share-modal"
      :files="mockFileArray"
    />
  </div>
</template>

<script>
import pdfjs from 'pdfjs-dist/webpack'
import throttle from 'lodash/throttle'

import { PUT_FILE_FOLDER, REMOVE_FILE_FOLDER } from '@/store/actions.type'
import { mapGetters } from 'vuex'

import PDFPage from '@/components/pdf-page.vue'
import FileListBatch from '@/components/file-list-batch'
import FolderModal from '@/components/modals/folder-modal'
import ShareModal from '@/components/modals/share-modal'

export default {
  name: 'pdf-doc',
  components: {
    PDFPage,
    FileListBatch,
    ShareModal,
    FolderModal
  },
  props: {
    fileID: String,
    acr: String
  },
  data () {
    return {
      pdf: null,
      pages: [],
      scale: 3,
      scrollTop: 0,
      clientHeight: 0
    }
  },
  computed: {
    url () {
      return `/api/file/${this.fileID}/download`
    },
    // TODO: File-List-Batch should be a File-List-Actions (my bad)
    // TODO: Get the actual file information for File-List-Batch.  I'd like to
    //  replace this with a file control component.
    mockFileArray () {
      return [{
        id: this.fileID,
        name: this.fileID,
        own: '-',
        date: { upload: '-' }
      }]
    },
    ...mapGetters(['activeGroup'])
  },
  methods: {
    increaseScale () {
      this.scale += 0.3
      // console.log('pdf-doc: scale: ', this.scale)
    },
    decreaseScale () {
      this.scale = Math.max(this.scale - 0.3, 0.1)
      // console.log('pdf-doc: scale: ', this.scale)
    },
    updateScrollBounds () {
      const { scrollTop, clientHeight } = this.$refs['pdfdoc']
      this.scrollTop = scrollTop
      this.clientHeight = clientHeight
      // console.log('pdf-doc: updated scroll bounds:', this.scrollTop, this.clientHeight)
    },
    pageInput (event) {
      const pageNumber = parseInt(event.target.value, 10) - 1
      if (isNaN(pageNumber) || pageNumber > this.pages.length || pageNumber < 0) return
      const pageOffset = this.$refs.pages[pageNumber].elementTop
      this.$refs['pdfdoc'].scrollTop = pageOffset
    },

    // file actions
    adjustFavorite (add) {
      return this.$store.dispatch(add ? PUT_FILE_FOLDER : REMOVE_FILE_FOLDER,
        { fid: this.fileID, name: '_favorites_', group: (this.gid || undefined) })
    },
    showShareModal () {
      this.$refs['shareModal'].show()
    },
    showFolderModal () {
      this.$refs['folderModal'].show()
    },
    createFolder (name) {
      this.$store.dispatch(PUT_FILE_FOLDER, { fid: this.fileID, name, group: this.activeGroup ? this.activeGroup.id : undefined })
    }
  },
  watch: {
    pdf (pdf) {
      this.pages = []
      let promises = []
      for (let i = 1; i <= pdf.numPages; i++) {
        promises.push(pdf.getPage(i))
      }
      Promise.all(promises).then(pages => {
        this.pages = pages
      }).catch(() => {
        // console.log(err)
      })
    },
    fileID () {
      pdfjs.getDocument(this.url).promise
        .then(pdf => {
          this.pdf = pdf
        })
        .catch(() => {
          // console.log(err)
        })
    }
  },
  created () {
    pdfjs.getDocument(this.url).promise
      .then(pdf => {
        this.pdf = pdf
      })
      .catch(() => {
        // console.log(err)
      })
  },
  mounted () {
    this.updateScrollBounds()
    this.throttledCallback = throttle(this.updateScrollBounds, 300)
    this.$el.addEventListener('scroll', this.throttledCallback, true)
    window.addEventListener('resize', this.throttledCallback, true)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.throttledCallback, true)
  }
}
</script>

<style lang="scss" scoped>

.pdf-doc {
  position: absolute;
  overflow: auto;
  width: 100%;
  top: 35px;
  bottom: 0;
  right: 0;
  left: 0;
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

svg {
  width: 100%;
  margin-top: -10px;
  height: 20px;
}

input {
  width: 50%;
}

</style>
