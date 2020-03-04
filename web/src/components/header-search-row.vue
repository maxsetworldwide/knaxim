<template>
  <b-row no-gutters class="w-90">
    <b-col>
      <b-row align-v="center">
        <b-col cols="1">
          <file-icon :folder="false" :webpage="webpage" :extention="ext" />
        </b-col>
        <b-col offset="1" cols="10" class="text-center ellipsis d-md-none">
          <b-link :to="`/file/` + id">{{ name }}</b-link>
        </b-col>
        <b-col class="d-none d-md-flex">
          <b-link :to="`/file/` + id">{{ name }}</b-link>
        </b-col>
      </b-row>
      <b-spinner v-if="matches.loading"></b-spinner>
      <ol v-else class="pl-3">
        <li v-for="(sum, indx) in rows" :key="indx">
          <b-row>
            <b-col class="line-no">
              <span>.{{ sum.lineNo + 1 }}</span>
            </b-col>
            <b-col class="pl-2">
              <span v-html="highlight(sum.text)"></span>
            </b-col>
          </b-row>
        </li>
      </ol>
    </b-col>
    <!-- TODO: Add the expand, collapse triagle at the bottom.
    <b-row align-v="center" class="text-center w-100">
      <b-col>
        <svg class="expand">
          <use href="../assets/app.svg#expand-tri"/>
        </svg>
      </b-col>
    </b-row>
    -->
  </b-row>
</template>
<script>
import fileIcon from '@/components/file-icon'
import { mapGetters } from 'vuex'

export default {
  name: 'header-search-row',
  components: {
    fileIcon
  },
  props: {
    webpage: {
      type: Boolean,
      default: false
    },
    name: {
      type: String,
      required: true
    },
    ext: {
      type: String,
      required: true
    },
    id: {
      type: String,
      required: true
    },
    find: {
      type: String,
      required: true
    },
    acr: {
      type: String,
      default: ''
    }
  },
  computed: {
    rows () {
      return this.matches.matched.slice(0, 4).map(row => {
        return {
          text: row.Content[0] || '',
          lineNo: row.Position
        }
      })
    },
    matches () {
      return this.allmatches[this.id]
    },
    ...mapGetters({
      allmatches: 'searchLines'
    })
  },
  methods: {
    highlight (summary) {
      summary = this.escapeSummary(summary)
      // Highlight an acronym and its subject OR every word; Largest Words First
      const pattern = this.acr
        ? `${this.find}|${this.acr}`
        : (match => {
          return match
            .split('"')
            .reduce((acc, phrase, indx) => {
              if (indx % 2 === 0) {
                return acc.concat(phrase.split(' '))
              } else {
                acc.push(phrase)
                return acc
              }
            }, [])
            .filter(word => word.length > 0)
            .sort((a, b) => b.length - a.length)
            .join('|')
        })(this.find)
      // this.find.split(' ').sort((a, b) => b.length - a.length).join('|')
      return summary.replace(new RegExp(pattern, 'gi'), match => {
        return `<span class="lite">${match}</span>`
      })
    },
    escapeSummary (summary) {
      // When adding more replacement keys, be sure to escape any special regex
      // characters to be replaced. e.g. adding a replacement for '?' should be
      // inserted as '\\?' as a key in the object (\ must be escaped in the
      // string so it reaches the regex)
      const replacements = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#x27;',
        '/': '&#x2F;'
      }
      const regex = new RegExp(Object.keys(replacements).join('|'), 'gi')
      summary = summary.replace(regex, match => {
        return replacements[match]
      })
      return summary
    }
  }
}
</script>

<style lang="scss" scoped>
.line-no {
  max-width: 3em;
  text-align: right;
  padding-right: 0px;
  direction: rtl;
  overflow: hidden;
}
ol {
  list-style: none;
}
li {
  line-height: 1.2em;
  padding-top: 0.6em;
}
.expand {
  width: 25px;
  height: 25px;
}
.lite {
  background: lightyellow;
}
.unrecognized {
  text-align: left;
}
.ellipsis {
  display: inline-block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
