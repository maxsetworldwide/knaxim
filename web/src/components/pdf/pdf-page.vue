<!--
pdf-page: an individual page for the pdf viewer

props:
  page: the page object given by pdfjs
  scale: the zoom level of the document
  scrollTop: the pdf viewer's current scrollTop
  clientHeight: the pdf viewer's current height (paired with scrollTop
                for determining visibility)
  sentenceHighlight: toggle the highlighting of sentences

events:
  'visible', pageNum: emitted when this page becomes visible (1 indexed)
  'matches', { pageNum, matches }: emitted upon finishing the search for matches

-->
<template>
  <div>
    <canvas :ref="canvasID" v-bind="canvasAttrs" />
    <!-- <div :style="textLayerDimStyle" class="text-layer" :ref="textLayerID" /> -->
    <pdf-text-layer
      @rendered="refreshTextLayer = false"
      v-bind="{
        sentenceHighlight,
        page,
        textLayerDimStyle,
        scale,
        refreshTextLayer
      }"
    />
  </div>
</template>

<script>
import PdfTextLayer from '@/components/pdf/pdf-text-layer'

export default {
  name: 'pdf-page',
  components: {
    PdfTextLayer
  },
  props: {
    page: Object,
    scale: Number,
    scrollTop: {
      type: Number,
      default: 0
    },
    clientHeight: {
      type: Number,
      default: 0
    },
    sentenceHighlight: {
      type: Boolean,
      default: true
    }
  },
  data () {
    return {
      elementTop: 0,
      elementHeight: 0,
      canvasOffsetLeft: 0,
      canvasOffsetTop: 0,
      canvas: null,
      refreshTextLayer: false
    }
  },
  computed: {
    isElementVisible () {
      const { elementTop, elementBottom, scrollTop, scrollBottom } = this
      if (!elementBottom) return
      return elementTop < scrollBottom && elementBottom > scrollTop
    },
    isElementFocused () {
      const {
        elementTop,
        elementBottom,
        scrollTop,
        scrollBottom,
        clientHeight
      } = this
      if (!elementBottom) return
      return (
        (elementTop > scrollTop && elementTop < scrollTop + clientHeight / 2) ||
        (elementBottom < scrollBottom &&
          elementBottom > scrollBottom - clientHeight / 2) ||
        (elementTop <= scrollTop && elementBottom >= scrollBottom)
      )
    },
    elementBottom () {
      return this.elementTop + this.elementHeight
    },
    scrollBottom () {
      return this.scrollTop + this.clientHeight
    },
    canvasID () {
      return 'canvas-' + this.page.pageNumber
    },
    canvasAttrs () {
      let { width, height } = this.viewport
      const style = this.canvasStyle
      return {
        width,
        height,
        style,
        class: 'pdf-page'
      }
    },
    canvasStyle () {
      const {
        width: actualSizeWidth,
        height: actualSizeHeight
      } = this.actualSizeViewport
      const pixelWidth = actualSizeWidth / this.pixelRatio
      const pixelHeight = actualSizeHeight / this.pixelRatio
      return `width: ${pixelWidth}px; height:${pixelHeight}px`
    },
    textLayerDimStyle () {
      const { canvasOffsetTop, canvasOffsetLeft } = this
      const viewport = this.actualSizeViewport
      const height = viewport.height
      const width = viewport.width

      const result = {
        left: canvasOffsetLeft + 'px',
        top: canvasOffsetTop + 'px',
        height: height + 'px',
        width: width + 'px'
      }
      return result
    },
    actualSizeViewport () {
      return this.viewport.clone({ scale: this.scale })
    },
    pixelRatio () {
      return window.devicePixelRatio || 1
    },
    sentenceStyle () {
      return this.sentenceHighlight ? 'sentenceOn' : 'sentenceOff'
    }
  },
  methods: {
    updateElementBounds () {
      const { offsetTop, offsetHeight } = this.$el
      this.elementTop = offsetTop
      this.elementHeight = offsetHeight
    },
    drawPage () {
      if (this.renderTask) return

      const viewport = this.viewport
      const canvasContext = this.$refs[this.canvasID].getContext('2d')
      const renderContext = { canvasContext, viewport }

      this.renderTask = this.page.render(renderContext)
      this.renderTask.promise.catch(() => {
        // console.log('pdf-page: renderTask failed: ', err)
        this.destroyRenderTask()
      })
    },
    destroyPage (page) {
      if (!page) return

      page._destroy()
      if (this.renderTask) {
        this.renderTask.cancel()
      }
    },
    destroyRenderTask () {
      if (!this.renderTask) return

      this.renderTask.cancel()
      delete this.renderTask
    }
  },
  watch: {
    isElementVisible (val) {
      if (val) {
        this.drawPage()
      }
    },
    isElementFocused (val) {
      if (val) {
        this.$emit('visible', this.page.pageIndex + 1)
      }
    },
    scale () {
      this.updateElementBounds()
      this.refreshTextLayer = true
    },
    scrollTop: 'updateElementBounds',
    clientHeight: 'updateElementBounds',
    page (newPage, oldPage) {
      this.destroyPage(oldPage)
    }
  },
  beforeDestroy () {
    this.destroyPage(this.page)
  },
  created () {
    this.viewport = this.page.getViewport({
      scale: this.scale / this.pixelRatio
    })
  },
  mounted () {
    this.canvas = this.$refs[this.canvasID]
    this.updateElementBounds()
    const { offsetTop, offsetLeft } = this.canvas
    this.canvasOffsetTop = offsetTop
    this.canvasOffsetLeft = offsetLeft
  },
  updated () {
    const { canvas } = this
    if (!canvas) return
    const { offsetTop, offsetLeft } = canvas
    this.canvasOffsetTop = offsetTop
    this.canvasOffsetLeft = offsetLeft
  }
}
</script>

<style lang="scss" scoped>
.pdf-page {
  display: block;
  margin: 0 auto;
}
</style>
