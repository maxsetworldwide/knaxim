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
