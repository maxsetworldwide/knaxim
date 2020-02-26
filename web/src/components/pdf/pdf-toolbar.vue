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
  'fit-height': fit to height button was pressed
  'fit-width': fit to width button was pressed
  'page-input', pageNumber: a page number was input and has been confirmed to
                  be a valid page input

-->
<template>
  <b-row align-v="end">
    <b-col class="d-none d-md-flex" offset="1" cols="1">
      <b-button @click="increaseScale">
        <svg>
          <use href="@/assets/app.svg#zoom-in"></use>
        </svg>
      </b-button>
    </b-col>
    <b-col class="d-none d-md-flex" cols="1">
      <b-button @click="decreaseScale">
        <svg>
          <use href="@/assets/app.svg#zoom-out"></use>
        </svg>
      </b-button>
    </b-col>
    <b-col class="d-none d-md-flex" cols="1">
      <b-button @click="fitWidth">
        <b-icon-arrow-left-right scale="1.4" />
      </b-button>
    </b-col>
    <b-col class="d-none d-md-flex" cols="1">
      <b-button @click="fitHeight">
        <b-icon-arrow-up-down scale="1.4" />
      </b-button>
    </b-col>
    <b-col offset="4" offset-md="0" cols="6" md="4">
      <h4 class="title text-center">{{ file.name }}</h4>
    </b-col>
    <b-col class="d-none d-md-flex" cols="2">
      <input
        :value="currPage"
        @input="onPageInput"
        min="1"
        :max="maxPages"
        type="number"
      />
      <span> / {{ maxPages }}</span>
    </b-col>

    <b-col cols="2" md="1">
      <file-actions singleFile :checkedFiles="[file]" />
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
    ...mapGetters(['activeGroup', 'getFolder'])
  },
  methods: {
    increaseScale () {
      this.$emit('scale-increase')
    },
    decreaseScale () {
      this.$emit('scale-decrease')
    },
    fitHeight () {
      this.$emit('fit-height')
    },
    fitWidth () {
      this.$emit('fit-width')
    },
    onPageInput (event) {
      const pageNumber = parseInt(event.target.value, 10)
      if (isNaN(pageNumber) || pageNumber < 1) return
      this.$emit('page-input', pageNumber)
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
}
</style>
