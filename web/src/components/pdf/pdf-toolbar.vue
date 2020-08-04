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
pdf-toolbar: provide actions and options for the pdf viewer

props:
  currPage: the current focused page. Use this to update the page selector
            with the current page.
  maxPages: the number of pages in the document
  id: the file ID of the document

events:
  'scale-increase': scale increase button was pressed
  'scale-decrease': scale decrease button was pressed
  'page-input', pageNumber: a page number was input and has been confirmed to
                  be a valid page input

-->
<template>
  <b-row align-v="end">
    <b-col class="d-none d-md-flex" md="1">
      <b-button v-if="navigation.prev" @click="openPrev" class="min-width-1em">
        <b-icon icon="arrow-bar-left" class="icon" />
      </b-button>
    </b-col>
    <b-col cols="2">
      <file-actions singleFile :checkedFiles="[file]" />
    </b-col>
    <b-col class="d-none d-md-flex" md="2">
      <input
        :value="currPage"
        @input="onPageInput"
        min="1"
        :max="maxPages"
        type="number"
        class="min-width-3em"
      />
      <span> / {{ maxPages }}</span>
    </b-col>
    <b-col offset="1" offset-md="0" cols="6" md="4">
      <h4 class="title text-center">{{ file.name }}</h4>
    </b-col>
    <b-col class="d-none d-md-flex" md="1">
      <b-button @click="increaseScale" class="min-width-3em">
        <svg>
          <use href="@/assets/app.svg#zoom-in"></use>
        </svg>
      </b-button>
    </b-col>
    <b-col class="d-none d-md-flex" md="1">
      <b-button @click="decreaseScale" class="min-width-3em">
        <svg>
          <use href="@/assets/app.svg#zoom-out"></use>
        </svg>
      </b-button>
    </b-col>
    <b-col class="d-none d-md-flex" md="1">
      <b-button v-if="navigation.next" @click="openNext" class="min-width-1em">
        <b-icon icon="arrow-bar-right" class="icon" />
      </b-button>
    </b-col>
  </b-row>
</template>

<script>
import FileActions from '@/components/file-actions'
import { mapGetters } from 'vuex'

export default {
  name: 'pdf-toolbar',
  components: {
    FileActions
  },
  props: {
    currPage: Number,
    maxPages: Number,
    file: Object
  },
  data () {
    return {
      pageInput: this.currPage
    }
  },
  computed: {
    navigation () {
      let nav = {
        prev: null,
        next: null
      }
      let indx = this.activeFiles.findIndex(af => af === this.file.id)
      if (indx > 0) {
        nav.prev = this.activeFiles[indx - 1]
      }
      if (indx >= 0 && indx < this.activeFiles.length - 1) {
        nav.next = this.activeFiles[indx + 1]
      }
      return nav
    },
    ...mapGetters(['activeFiles'])
  },
  methods: {
    increaseScale () {
      this.$emit('scale-increase')
    },
    decreaseScale () {
      this.$emit('scale-decrease')
    },
    onPageInput (event) {
      const pageNumber = parseInt(event.target.value, 10)
      if (isNaN(pageNumber) || pageNumber < 1) return
      this.$emit('page-input', pageNumber)
    },
    openPrev () {
      if (this.navigation.prev) {
        this.$router.push(`/file/${this.navigation.prev}`)
      }
    },
    openNext () {
      if (this.navigation.next) {
        this.$router.push(`/file/${this.navigation.next}`)
      }
    }
  },
  watch: {
    currPage (val) {
      this.pageInput = val
    }
  }
}
</script>

<style scoped lang="scss">
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

.title {
  width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  position: absolute;
  bottom: 0;
}

.min-width-3em {
  min-width: 3em;
}
</style>
