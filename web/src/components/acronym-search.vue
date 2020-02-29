<template>
  <div v-if="result.length">
    <slot :result="result">
    </slot>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { ACRONYMS } from '@/store/actions.type'

export default {
  name: 'acronym-search',
  props: {
    phrase: {
      type: String,
      default: '',
      required: true
    }
  },
  watch: {
    phrase: function (n, o) {
      if (n !== o) {
        this.search({ acronym: n })
      }
    }
  },
  methods: {
    ...mapActions({
      search: ACRONYMS
    })
  },
  computed: {
    result () {
      return this.acronymResults
    },
    ...mapGetters(['acronymResults'])
  }
}
</script>
