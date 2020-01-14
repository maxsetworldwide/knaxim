<template>
  <div>
    <canvas :ref="canvasID" v-bind="canvasAttrs"/>
    <div :style="textLayerDimStyle" class="text-layer" :ref="textLayerID"/>
  </div>
</template>

<script>
import pdfjs from 'pdfjs-dist/webpack'
import { mapGetters } from 'vuex'

export default {
  name: 'pdf-page',
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
    }
  },
  data () {
    return {
      elementTop: 0,
      elementHeight: 0,
      renderTextLayerTask: null,
      textContent: null,
      textLayerDimStyle: {},
      textSpans: [],
      textContentItemsStr: [],
      staleTextLayer: true
    }
  },
  computed: {
    isElementVisible () {
      const { elementTop, elementBottom, scrollTop, scrollBottom } = this
      if (!elementBottom) return
      return elementTop < scrollBottom && elementBottom > scrollTop
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
    textLayerID () {
      return 'text-container-' + this.page.pageNumber
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
    actualSizeViewport () {
      return this.viewport.clone({ scale: this.scale })
    },
    pixelRatio () {
      return window.devicePixelRatio || 1
    },
    ...mapGetters(['currentSearch'])
  },
  methods: {
    updateElementBounds () {
      const { offsetTop, offsetHeight } = this.$el
      this.elementTop = offsetTop
      this.elementHeight = offsetHeight
      // console.log('updated element bounds:', this.elementTop, this.elementHeight)
    },
    setDimStyle () {
      // TODO: make computed, problem is $refs. add watcher?
      const canvas = this.$refs[this.canvasID]
      const viewport = this.actualSizeViewport
      const height = viewport.height
      const width = viewport.width

      this.textLayerDimStyle = {
        left: canvas.offsetLeft + 'px',
        top: canvas.offsetTop + 'px',
        height: height + 'px',
        width: width + 'px'
      }
    },
    renderText () {
      if (this.textContent) {
        this.setDimStyle()
        this.$refs[this.textLayerID].innerHTML = ''
        this.textSpans = []
        pdfjs.renderTextLayerTask = pdfjs.renderTextLayer({
          textContent: this.textContent,
          viewport: this.page.getViewport({ scale: this.scale / this.pixelRatio }),
          container: this.$refs[this.textLayerID],
          textDivs: this.textSpans,
          textContentItemsStr: this.textContentItemsStr
        })
        this.highlightMatches()
      }
    },
    highlightMatches () {
      const search = this.currentSearch.trim().toLowerCase()
      const linkObj = {
        'DD Form 460': 'dd0460.pdf',
        'DD Form 0499': 'dd0499.pdf',
        'DD Form 0294': 'dd0294.pdf',
        'AF Form 4080': 'af4080.pdf',
        'AF Form 2407': 'af2407.pdf',
        'AF Form 538': 'af538.pdf',
        'AF Form 228': 'af228.pdf',
        'AF Form 35': 'af35.pdf'
      }
      const prefix = 'file://home/demo/Documents/'
      let linkObjHasTerms = false
      for (let key in linkObj) {
        if (linkObj.hasOwnProperty(key)) {
          linkObjHasTerms = true
          break
        }
      }
      if (search.length === 0 && !linkObjHasTerms) return
      const { textSpans, textContentItemsStr: contentArr } = this
      for (let i = 0; i < textSpans.length; i++) {
        let span = textSpans[i]
        let content = contentArr[i]
        let currOffset = 0
        let nextMatchOffset
        do {
          if (search.length > 0) {
            nextMatchOffset = content.toLowerCase().indexOf(search, currOffset)
          } else {
            nextMatchOffset = -1
          }
          let nextKey = search
          let link = false
          for (let key in linkObj) {
            if (linkObj.hasOwnProperty(key)) {
              let candidateOffset = content.indexOf(key, currOffset)
              if (nextMatchOffset === -1 || (candidateOffset < nextMatchOffset && candidateOffset !== -1)) {
                nextMatchOffset = candidateOffset
                nextKey = key
                link = true
              }
            }
          }
          if (nextMatchOffset > -1) {
            if (currOffset === 0) {
              span.textContent = ''
            }
            if (link) {
              this.appendTextChild(i, currOffset, nextMatchOffset, '')
              this.appendTextChild(i, nextMatchOffset, nextMatchOffset + nextKey.length, 'link', prefix + linkObj[nextKey])
            } else {
              this.appendTextChild(i, currOffset, nextMatchOffset, '')
              this.appendTextChild(i, nextMatchOffset, nextMatchOffset + nextKey.length, 'match')
            }
            currOffset = nextMatchOffset + nextKey.length
          }
        } while (nextMatchOffset > -1)
        if (currOffset > 0) {
          this.appendTextChild(i, currOffset, content.length, '')
        }
      }
    },
    appendTextChild (spanIdx, from, to, className, link) {
      let parentSpan = this.textSpans[spanIdx]
      let parentContent = this.textContentItemsStr[spanIdx]
      let substring = parentContent.substring(from, to)
      let textNode = document.createTextNode(substring)
      let childNode = document.createElement('span')
      if (className && !link) {
        childNode.className = className
        childNode.appendChild(textNode)
        parentSpan.appendChild(childNode)
      } else if (link) {
        childNode.className = className
        let linkNode = document.createElement('a')
        linkNode.setAttribute('href', link)
        linkNode.appendChild(textNode)
        childNode.appendChild(linkNode)
        parentSpan.appendChild(childNode)
      } else {
        parentSpan.appendChild(textNode)
      }
    },
    drawPage () {
      if (this.renderTask) return

      const viewport = this.viewport
      const canvasContext = this.$refs[this.canvasID].getContext('2d')
      const renderContext = { canvasContext, viewport }

      this.renderTask = this.page.render(renderContext)
      this.renderTask.promise
        .then(() => {
          return this.page.getTextContent()
        })
        .then((content) => {
          this.textContent = content
          this.renderText()
        })
        .catch(() => {
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
    scale () {
      this.updateElementBounds()
      if (this.isElementVisible) {
        this.staleTextLayer = true
      }
    },
    scrollTop: 'updateElementBounds',
    clientHeight: 'updateElementBounds',
    page (newPage, oldPage) {
      this.destroyPage(oldPage)
    }
  },
  updated () {
    // text layer is rendered after component update due to requiring data
    // from non-reactive canvas element in $refs
    if (this.staleTextLayer) {
      this.staleTextLayer = false
      this.renderText()
    }
  },
  beforeDestroy () {
    this.destroyPage(this.page)
  },
  created () {
    this.viewport = this.page.getViewport({ scale: this.scale / this.pixelRatio })
  },
  mounted () {
    this.updateElementBounds()
  }
}
</script>

<style lang="scss" scoped>

.text-layer {
   position: absolute;
   display: inline-block;
   left: 0;
   top: 0;
   right: 0;
   bottom: 0;
   overflow: hidden;
   opacity: 0.5;
   line-height: 1.0;
}

.pdf-page {
  display: block;
  margin: 0 auto;
}

</style>

<style>

.text-layer > span {
  color: transparent;
  position: absolute;
  white-space: pre;
  cursor: text;
  -webkit-transform-origin: 0% 0%;
  -moz-transform-origin: 0% 0%;
  -o-transform-origin: 0% 0%;
  -ms-transform-origin: 0% 0%;
  transform-origin: 0% 0%;
}

.text-layer ::selection {
  background: rgb(0, 0, 255);
}

.match {
  background-color: #75ADCB;
}

.link {
  background-color: green;
  cursor: pointer;
}

</style>
