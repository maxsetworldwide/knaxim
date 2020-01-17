<template>
  <div class="h-100">
    <b-container class="h-100" fluid>
      <b-row>
        <b-col offset-md="5" cols="1">
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

        <b-col offset-md="2" cols="1">
          <file-list-batch fileSelected singleFile
            :checkedFiles="mockFileArray"
            :removeFavorite="isFavorite"
            @favorite="adjustFavorite"
            @add-folder="showFolderModal"
            @share-file="showShareModal"
          />
        </b-col>
      </b-row>
      <b-row class="h-100 pdf-row" align-h="start">
        <b-col cols="2">
          <PDFResultList class="h-100 result-list" :matchList="matchList" @select="onMatchSelect"/>
        </b-col>
        <b-col>
          <div class="pdf-doc" :ref="'pdfdoc'">
            <PDFPage
              v-for="page in pages"
              v-bind="{ page, scale, scrollTop, clientHeight }"
              :key="page.pageNumber"
              ref="pages"
              @matches="handleMatches"
            />
          </div>
        </b-col>
      </b-row>
    </b-container>
    <folder-modal
      ref="folderModal"
      id="newFolderModal"
      @new-folder="createFolder"
    />
    <share-modal
      hideList
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
import PDFResultList from './pdf-result-list'
import FileListBatch from '@/components/file-list-batch'
import FolderModal from '@/components/modals/folder-modal'
import ShareModal from '@/components/modals/share-modal'

export default {
  name: 'pdf-doc',
  components: {
    PDFPage,
    PDFResultList,
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
      scale: 1.25,
      scrollTop: 0,
      clientHeight: 0,
      matchContexts: {}
    }
  },
  computed: {
    isFavorite () {
      let favorite = (this.$store.state.folder.user['_favorites_'] || [])
      return favorite.reduce((acc, id) => {
        return acc || id === this.fileID
      }, false)
    },
    url () {
      return `/api/file/${this.fileID}/download`
    },
    matchList () {
      let result = []
      for (let page in this.matchContexts) {
        let matchList = this.matchContexts[page]
        matchList.forEach(match => {
          match.page = page
          result.push(match)
        })
      }
      return result
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
    },
    decreaseScale () {
      this.scale = Math.max(this.scale - 0.3, 0.1)
    },
    updateScrollBounds () {
      const { scrollTop, clientHeight } = this.$refs['pdfdoc']
      this.scrollTop = scrollTop
      this.clientHeight = clientHeight
    },
    pageInput (event) {
      const pageNumber = parseInt(event.target.value, 10) - 1
      if (isNaN(pageNumber) || pageNumber > this.pages.length || pageNumber < 0) return
      const pageOffset = this.$refs.pages[pageNumber].elementTop
      // console.log('page ref:', this.$refs.pages[pageNumber])
      this.$refs['pdfdoc'].scrollTop = pageOffset
    },
    onMatchSelect (match) {
      const offset = match.span.offsetTop + this.$refs.pages[match.page].elementTop
      this.$refs['pdfdoc'].scrollTop = offset - 20
    },
    handleMatches (payload) {
      this.$set(this.matchContexts, payload.pageNum, payload.matches)
      // console.log('matchContexts: ', this.matchContexts)
    },
    // file actions
    adjustFavorite (add) {
      return this.$store.dispatch(add ? PUT_FILE_FOLDER : REMOVE_FILE_FOLDER,
        { fid: this.fileID, name: '_favorites_' })
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
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
}

.pdf-row {
  margin-top: 25px;
}

.result-list {
  overflow: auto;
  position: absolute;
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
