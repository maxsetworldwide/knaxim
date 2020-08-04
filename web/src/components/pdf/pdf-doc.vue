<!--
// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
-->
<!--
  events:
    no-view: emitted when the request for the view from the server returns
             with an error
-->
<template>
  <div class="h-100">
    <b-container class="h-100" fluid>
      <pdf-toolbar
        class="pdf-toolbar"
        :currPage="currPage"
        :maxPages="pages.length"
        :file="getFile"
        @scale-increase="increaseScale"
        @scale-decrease="decreaseScale"
        @fit-height="fitToHeight"
        @fit-width="fitToWidth"
        @page-input="pageInput"
      />
      <b-row class="pdf-row" align-h="start">
        <b-col
          v-if="currentSearch.length > 0"
          class="d-none d-lg-inline"
          cols="2"
        >
          <pdf-result-list
            class="h-100 result-list"
            :matchList="matchList"
            @select="onMatchSelect"
            @highlight="sentenceHighlight = $event"
          />
        </b-col>
        <b-col>
          <div class="pdf-doc" :ref="'pdfdoc'">
            <pdf-page
              class="pdf-page"
              v-for="page in pages"
              v-bind="{
                page,
                scale,
                scrollTop,
                clientHeight,
                sentenceHighlight
              }"
              :key="page.pageNumber"
              ref="pages"
              @matches="handleMatches"
              @visible="handleVisible"
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
import { GET_FILE } from '@/store/actions.type'
import { mapActions, mapGetters } from 'vuex'

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
      matchContexts: {},
      sentenceHighlight: true
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
        matchList.forEach((match) => {
          match.page = page
          result.push(match)
        })
      }
      return result
    },
    getFile () {
      this.fetchFile({ id: this.fileID })
      const file = this.populateFiles(this.fileID)
      if (!file) return {}
      return file
    },
    ...mapGetters(['populateFiles', 'currentSearch'])
  },
  methods: {
    increaseScale () {
      this.scale += 0.3
    },
    decreaseScale () {
      this.scale = Math.max(this.scale - 0.3, 0.1)
    },
    fitToWidth () {
      const pixelRatio = window.devicePixelRatio || 1
      this.scale = (this.$el.clientWidth * pixelRatio) / 800
    },
    fitToHeight () {
      const pixelRatio = window.devicePixelRatio || 1
      this.scale = (this.$el.clientHeight * pixelRatio) / 800
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
    fetchPdf () {
      pdfjs
        .getDocument(this.url)
        .promise.then((pdf) => {
          this.pdf = pdf
        })
        .catch(() => {
          this.$emit('no-view')
          // console.log(err)
        })
    },
    ...mapActions({
      fetchFile: GET_FILE
    })
  },
  watch: {
    pdf (pdf) {
      this.pages = []
      let promises = []
      for (let i = 1; i <= pdf.numPages; i++) {
        promises.push(pdf.getPage(i))
      }
      Promise.all(promises)
        .then((pages) => {
          this.pages = pages
          this.fitToWidth()
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
.pdf-toolbar {
  height: 5%;
}

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
  height: 90%;
}

.result-list {
  overflow: auto;
  position: absolute;
}
</style>
