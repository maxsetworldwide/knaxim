<!--
  events:
    no-view: emitted when the request for the view from the server returns
             with an error
-->
<template>
  <div class="h-100">
    <b-container class="h-100" fluid>
      <pdf-toolbar
        :currPage="currPage"
        :maxPages="pages.length"
        :id="fileID"
        :name="name"
        @scale-increase="increaseScale"
        @scale-decrease="decreaseScale"
        @page-input="pageInput"
      />
      <b-row class="h-100 pdf-row" align-h="start">
        <b-col cols="2">
          <pdf-result-list
            class="h-100 result-list"
            :matchList="matchList"
            @select="onMatchSelect"
          />
        </b-col>
        <b-col>
          <div class="pdf-doc" :ref="'pdfdoc'">
            <pdf-page
              class="pdf-page"
              v-for="page in pages"
              v-bind="{ page, scale, scrollTop, clientHeight }"
              :key="page.pageNumber"
              ref="pages"
              @matches="handleMatches"
              @visible="handleVisible"
              @viewport="handleViewport"
            />
          </div>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import pdfjs from 'pdfjs-dist/webpack'
import throttle from 'lodash/throttle'

import PdfPage from './pdf-page.vue'
import PdfResultList from './pdf-result-list'
import PdfToolbar from './pdf-toolbar'

import FileService from '@/service/file'

export default {
  name: 'pdf-doc',
  components: {
    PdfPage,
    PdfResultList,
    PdfToolbar
  },
  props: {
    fileID: String,
    acr: String
  },
  data () {
    return {
      pdf: null,
      name: '',
      pages: [],
      currPage: 1,
      scale: 3,
      scrollTop: 0,
      clientHeight: 0,
      matchContexts: {}
    }
  },
  computed: {
    url () {
      return FileService.viewURL({ fid: this.fileID })
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
    }
    // defaultViewport () {
    // if (!this.pages.length) {
    // return {
    // width: 0,
    // height: 0
    // }
    // }
    // let vp = this.pages[0].getViewport({ scale: this.scale })
    // return this.pages[0].getViewport({ scale: this.scale })
    // }
  },
  methods: {
    increaseScale () {
      this.scale += 0.3
    },
    decreaseScale () {
      this.scale = Math.max(this.scale - 0.3, 0.1)
    },
    pageHeightScale (vp) {
      let pixelRatio = window.devicePixelRatio || 1
      console.log('vp:', vp)
      return this.$el.clientWidth * pixelRatio * (1.5 / vp.width)
    },
    pageInput (pageNumber) {
      if (
        isNaN(pageNumber) ||
        pageNumber > this.pages.length ||
        pageNumber < 1
      ) {
        return
      }
      const pageOffset = this.$refs.pages[pageNumber - 1].elementTop
      // console.log('page ref:', this.$refs.pages[pageNumber])
      this.$refs['pdfdoc'].scrollTop = pageOffset
      this.currPage = pageNumber
    },
    updateScrollBounds () {
      const { scrollTop, clientHeight } = this.$refs['pdfdoc']
      this.scrollTop = scrollTop
      this.clientHeight = clientHeight
    },
    onMatchSelect (match) {
      const offset =
        match.span.offsetTop + this.$refs.pages[match.page].elementTop
      this.$refs['pdfdoc'].scrollTop = offset - 20
    },
    handleMatches (payload) {
      this.$set(this.matchContexts, payload.pageNum, payload.matches)
      // console.log('matchContexts: ', this.matchContexts)
    },
    handleVisible (num) {
      this.currPage = num
    },
    handleViewport (vp) {
      this.scale = this.pageHeightScale(vp)
    },
    fetchPdf () {
      FileService.info({ fid: this.fileID }).then(({ data }) => {
        this.name = data.file.name
      })
      pdfjs
        .getDocument(this.url)
        .promise.then(pdf => {
          this.pdf = pdf
        })
        .catch(() => {
          this.$emit('no-view')
          // console.log(err)
        })
    }
  },
  watch: {
    pdf (pdf) {
      this.pages = []
      let promises = []
      for (let i = 1; i <= pdf.numPages; i++) {
        promises.push(pdf.getPage(i))
      }
      Promise.all(promises)
        .then(pages => {
          this.pages = pages
        })
        .catch(() => {
          // console.log(err)
        })
    },
    fileID () {
      this.fetchPdf()
    }
  },
  created () {
    this.fetchPdf()
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

.pdf-page {
  margin-bottom: 10px;
}

.pdf-row {
  margin-top: 25px;
}

.result-list {
  overflow: auto;
  position: absolute;
}
</style>
