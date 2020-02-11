<template>
  <b-container fluid class="header-search-list">
    <b-row no-gutters class="w-90"
        v-for="item in rows"
        :key="item.id">
      <b-col cols="1" class="">
        <file-icon :folder="false" :webpage="item.webpage" :extention="item.ext"/>
      </b-col>
      <b-col cols="11">
        <b-link :to="`/file/` + item.id">{{ item.name }}</b-link>
        <ol>
          <li v-for="(sum, indx) in item.summary" :key="indx" v-html="highlight(sum)"></li>
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
  </b-container>
</template>

<script>
import fileIcon from '@/components/file-icon'
import { FILES_SEARCH } from '@/store/actions.type'
import { mapGetters } from 'vuex'

export default {
  components: {
    fileIcon
  },
  props: {
    find: {
      type: String,
      required: true
    },
    acr: String
  },
  name: 'header-search-list',
  // TODO: Merge the watch and before Mount logic or create a method to handle
  //  both use cases.  One is for first render the other is for prop change.
  watch: {
    find (val) {
      this.$store.dispatch(FILES_SEARCH, { find: val, acr: this.acr })
        .catch(() => {
          // console.log(`Error: ${message}`)
        })
    },
    activeGroup () {
      this.$store.dispatch(FILES_SEARCH, { find: this.find, acr: this.acr })
    }
  },
  beforeMount () {
    this.$store.dispatch(FILES_SEARCH, { find: this.find, acr: this.acr })
      .catch(() => {
        // console.log(`Error: ${message}`)
      })
  },
  computed: {
    rows () {
      return this.searchMatches.map((file) => {
        let splits = file.name.split('.')
        let row = {
          id: file.id,
          name: file.name,
          webpage: file.name.includes('/'),
          ext: splits.length > 1 ? splits[splits.length - 1] : '',
          summary: ''
        }
        if (file.lines) {
          row.summary = file.lines.slice(0, 4).map(line => line.Content[0])
        }
        return row
      })
    },
    ...mapGetters(['searchMatches', 'activeGroup'])
  },
  methods: {
    highlight (summary) {
      summary = this.escapeSummary(summary)
      // Highlight an acronym and its subject OR every word; Largest Words First
      const pattern = this.acr ? `${this.find}|${this.acr}`
        : ((match) => {
          return match.split('"').reduce((acc, phrase, indx) => {
            if (indx % 2 === 0) {
              return acc.concat(phrase.split(' '))
            } else {
              acc.push(phrase)
              return acc
            }
          }, []).filter((word) => word.length > 0).sort((a, b) => b.length - a.length).join('|')
        })(this.find)
        // this.find.split(' ').sort((a, b) => b.length - a.length).join('|')
      return summary.replace(new RegExp(pattern, 'gi'), (match) => {
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
        '\'': '&#x27;',
        '/': '&#x2F;'
      }
      const regex = new RegExp(Object.keys(replacements).join('|'), 'gi')
      summary = summary.replace(regex, (match) => {
        return replacements[match]
      })
      return summary
    }
  }
}
</script>

<style lang="scss">
.header-search-list {
  li {
    line-height: 1.2em;
    padding-top: .3em;
  }
  .expand {
    width: 25px;
    height: 25px;
  }
  svg {
    width: 50px;
    height: 50px;
  }
  .lite {
    background: lightyellow
  }
  .unrecognized {
    text-align: left;
  }
}
</style>
