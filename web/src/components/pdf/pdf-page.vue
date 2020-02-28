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
    <div :style="textLayerDimStyle" class="text-layer" :ref="textLayerID" />
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
      renderTextLayerTask: null,
      textContent: null,
      textLayerDimStyle: {},
      textSpans: [],
      textContentItemsStr: [],
      joinedContent: '',
      matches: [],
      sentenceBounds: [],
      staleTextLayer: true
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
    sentenceStyle () {
      return this.sentenceHighlight ? 'sentenceOn' : 'sentenceOff'
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
    findMatches () {
      if (this.currentSearch.length === 0) return []
      const search = compileSearchTerms(this.currentSearch)
      this.textContentItemsStr = this.textContentItemsStr.map(str => {
        if (/\S/.test(str)) {
          return str
        } else {
          return ''
        }
      })
      this.joinedContent = this.textContentItemsStr.join('').toLowerCase()
      const { textSpans, joinedContent } = this
      this.sentenceBounds = findSentences()
      const sentenceBounds = this.sentenceBounds

      // take search query string and turn it into an array of terms, combining
      // terms with quotes.
      function compileSearchTerms (searchQuery) {
        const regex = /[^\s"]+|"([^"]*)"/g
        let result = []
        let match = regex.exec(searchQuery)
        while (match !== null) {
          result.push(match[1] ? match[1] : match[0])
          match = regex.exec(searchQuery)
        }
        result = result
          .filter(term => term.length > 0) // just in case
          .map(term => {
            return term.toLowerCase()
          })
        return result
      }

      function getSpanFromJump (
        globalDelta,
        globalStart,
        spanStart,
        localStart
      ) {
        const globalEnd = globalStart + globalDelta
        while (globalDelta > 0 && spanStart < textSpans.length) {
          let toNextSpan = textSpans[spanStart].innerText.length - localStart
          if (toNextSpan < globalDelta) {
            spanStart++
            localStart = 0
            globalDelta -= toNextSpan
          } else {
            localStart += globalDelta
            globalDelta = 0
          }
        }
        if (spanStart >= textSpans.length) {
          spanStart = textSpans.length - 1
        }
        return {
          span: spanStart,
          offset: localStart,
          global: globalEnd
        }
      }

      function findSentences () {
        const punct = ['.', '!', '?']
        let result = []
        let localOffset = 0
        let globalOffset = 0
        let span = 0
        let minOffset = -1
        do {
          minOffset = -1
          punct.forEach(p => {
            const candidate = joinedContent.indexOf(p, globalOffset)
            if (
              candidate !== -1 &&
              (candidate <= minOffset || minOffset === -1)
            ) {
              minOffset = candidate
            }
          })
          if (minOffset > -1) {
            const next = getSpanFromJump(
              minOffset - globalOffset,
              globalOffset,
              span,
              localOffset
            )
            result.push({
              start: {
                span: span,
                offset: localOffset,
                global: globalOffset
              },
              end: {
                span: next.span,
                offset: next.offset,
                global: next.global
              }
            })
            const nextStart = getSpanFromJump(
              1,
              next.global,
              next.span,
              next.offset
            )
            globalOffset = nextStart.global
            localOffset = nextStart.offset
            span = nextStart.span
          }
        } while (minOffset > 0)
        const offsetToEnd = joinedContent.length - globalOffset
        if (offsetToEnd > 0) {
          const final = getSpanFromJump(
            joinedContent.length - globalOffset,
            globalOffset,
            span,
            localOffset
          )
          result.push({
            start: {
              span: span,
              offset: localOffset,
              global: globalOffset
            },
            end: {
              span: final.span,
              offset: final.offset,
              global: final.global
            }
          })
        }
        return result
      }

      function getNextMatch (globalIdx, localIdx, spanIdx) {
        // using the current offsets, get the next search term offset
        // try to keep this precedural without side effects
        let minOffset = -1
        let nextTermLength = -1
        search.forEach(currTerm => {
          if (currTerm.length === 0) return
          let candidateOffset = joinedContent.indexOf(currTerm, globalIdx)
          if (
            (candidateOffset <= minOffset && candidateOffset !== -1) ||
            (minOffset === -1 && candidateOffset !== -1)
          ) {
            // be sure to favor larger word e.g. for searching 'to' and 'tomorrow'
            if (
              candidateOffset === minOffset &&
              nextTermLength > currTerm.length
            ) {
              return
            }
            minOffset = candidateOffset // global offset
            nextTermLength = currTerm.length // keyword offset
          }
        })
        if (minOffset === -1) {
          // not found
          return null
        }
        // find sentence bounds and span idx
        // span idx
        let start = getSpanFromJump(
          minOffset - globalIdx,
          globalIdx,
          spanIdx,
          localIdx
        )
        // get end point
        let end = getSpanFromJump(
          nextTermLength,
          start.global,
          start.span,
          start.offset
        )

        // sentence bounds
        let sentenceIdx = 0
        let sentenceFound = false
        for (
          ;
          sentenceIdx < sentenceBounds.length && !sentenceFound;
          sentenceIdx++
        ) {
          const curr = sentenceBounds[sentenceIdx]
          sentenceFound =
            curr.start.global <= start.global && curr.end.global >= end.global
        }
        if (!sentenceFound) {
          // console.log('pdf-page: sentence not found: start:', start, 'end:', end)
        }
        return {
          start,
          end,
          sentence: sentenceIdx - 1
        }
      }

      let spanIdx = 0
      let localIdx = 0
      let globalIdx = 0
      let matches = []
      let match = getNextMatch(globalIdx, localIdx, spanIdx)
      while (match) {
        matches.push(match)
        spanIdx = match.end.span
        localIdx = match.end.offset
        globalIdx = match.end.global
        match = getNextMatch(globalIdx, localIdx, spanIdx)
      }
      return matches
    },
    sendMatches () {
      let matchContexts = []
      for (const matchIdx in this.matches) {
        let match = this.matches[matchIdx]
        let sentence = this.sentenceBounds[match.sentence]
        let sentenceText = this.joinedContent.substring(
          sentence.start.global,
          sentence.end.global
        )
        let span = this.textSpans[this.matches[matchIdx].start.span]
        matchContexts.push({
          sentence: {
            text: sentenceText,
            start: sentence.start,
            end: sentence.end
          },
          match,
          span
        })
      }
      this.$emit('matches', {
        pageNum: this.page.pageIndex,
        matches: matchContexts
      })
    },
    renderText () {
      if (this.textContent) {
        this.staleTextLayer = false
        this.setDimStyle()
        this.$refs[this.textLayerID].innerHTML = ''
        this.textSpans = []
        this.textContentItemsStr = []
        pdfjs.renderTextLayerTask = pdfjs.renderTextLayer({
          textContent: this.textContent,
          viewport: this.page.getViewport({
            scale: this.scale / this.pixelRatio
          }),
          container: this.$refs[this.textLayerID],
          textDivs: this.textSpans,
          textContentItemsStr: this.textContentItemsStr
        })
        this.matches = this.findMatches()
        this.highlightMatches()
        this.sendMatches()
      }
    },
    highlightMatches () {
      const sentSet = this.matches.reduce((acc, curr) => {
        return acc.add(curr.sentence)
      }, new Set())

      // create sentence segments
      const sentences = [...sentSet].map(idx => {
        return this.sentenceBounds[idx]
      })
      let segments = {}
      sentences.forEach(s => {
        let from = s.start.offset
        let to
        for (let spanIdx = s.start.span; spanIdx <= s.end.span; spanIdx++) {
          if (spanIdx === s.end.span) {
            to = s.end.offset
          } else {
            to = this.textContentItemsStr[spanIdx].length
          }
          if (!segments[spanIdx]) {
            segments[spanIdx] = []
          }
          segments[spanIdx].push({
            start: from,
            end: to,
            type: this.sentenceStyle
            // the 'type' portion of this object will determine the CSS style
            // used in appendTextChild()
          })
          from = 0
        }
      })
      // create keyword segments
      this.matches.forEach(match => {
        let segList = segments[match.start.span]
        let newStartSegList = []
        let i = 0
        while (i < segList.length && match.start.offset > segList[i].end) {
          newStartSegList.push(segList[i])
          i++
        }
        if (i === segList.length) {
          // console.log('keyword not found in segment')
          return
        }
        // cut sentence segment
        let from = segList[i].start
        let to = match.start.offset
        if (from < to) {
          newStartSegList.push({
            start: from,
            end: to,
            type: segList[i].type
          })
        }
        // start keyword segment
        from = match.start.offset
        if (match.start.span === match.end.span) {
          to = match.end.offset
        } else {
          to = this.textContentItemsStr[match.start.span].length
        }
        newStartSegList.push({
          start: from,
          end: to,
          type: 'keyword'
        })
        // push segments until we reach the end of the match
        for (let i = match.start.span + 1; i < match.end.span; i++) {
          const interSeg = [
            {
              start: 0,
              end: this.textContentItemsStr[i].length,
              type: 'keyword'
            }
          ]
          segments[i] = interSeg
        }
        // push tail end of keyword if we crossed spans
        let newTailSeg = []
        if (match.start.span !== match.end.span) {
          newTailSeg.push({
            start: 0,
            end: match.end.offset,
            type: 'keyword'
          })
        }
        // push tail end of sentence
        // find segment in match.end with segment end > match end offset
        segList = segments[match.end.span]
        i = 0
        while (i < segList.length && match.end.offset > segList[i].end) {
          i++
        }
        if (i === segList.length) {
          // console.log('match end not found in segment')
          return
        }
        from = match.end.offset
        to = segList[i].end
        newTailSeg.push({
          start: from,
          end: to,
          type: segList[i].type
        })
        // push rest of the last segment
        i++
        while (i < segList.length) {
          newTailSeg.push(segList[i])
          i++
        }
        // apply changes
        if (match.start.span === match.end.span) {
          newTailSeg.forEach(seg => {
            newStartSegList.push(seg)
          })
          segments[match.start.span] = newStartSegList
        } else {
          segments[match.start.span] = newStartSegList
          segments[match.end.span] = newTailSeg
        }
      })
      // fill in segments that have no highlighting
      for (let spanIdx in segments) {
        let spanContent = this.textContentItemsStr[spanIdx]
        let segmentList = segments[spanIdx]
        let newList = []
        let from = 0
        let to
        segmentList.forEach(seg => {
          to = seg.start
          if (from !== to) {
            newList.push({
              start: from,
              end: to,
              type: ''
            })
          }
          newList.push(seg)
          from = seg.end
        })
        const lastElementEnd = newList.slice(-1)[0].end
        if (lastElementEnd !== spanContent.length) {
          newList.push({
            start: lastElementEnd,
            end: spanContent.length,
            type: ''
          })
        }
        segments[spanIdx] = newList
      }

      // segments all documented. can now edit spans solely based on the object
      for (let spanIdx in segments) {
        let span = this.textSpans[spanIdx]
        span.textContent = ''
        segments[spanIdx].forEach(seg => {
          this.appendTextChild(spanIdx, seg.start, seg.end, seg.type)
        })
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
      this.staleTextLayer = true
    },
    sentenceHighlight: 'renderText',
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
      this.renderText()
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
    this.updateElementBounds()
    this.page.getTextContent().then(content => {
      this.textContent = content
      this.renderText()
    })
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
  line-height: 1;
}

.pdf-page {
  display: block;
  margin: 0 auto;
}
</style>

<style lang="scss">
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

.keyword {
  background-color: goldenrod;
}

.sentenceOn {
  background-color: $app-clr2;
}

.sentenceOff {
  background-color: transparent;
}

.link {
  background-color: green;
  cursor: pointer;
}
</style>
