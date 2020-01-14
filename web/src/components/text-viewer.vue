<!--
text-viewer: a file viewer for reading the plain text version of a record

this component should be instantiated via the router, with the file ID as a
parameter

future work:
  Add watcher for route change if we will have linking between documents
    This is due to vue reusing the same component instance
    https://router.vuejs.org/guide/essentials/dynamic-matching.html#reacting-to-params-changes
-->

<template>
  <div>
    <b-container fluid class="header">
      <b-row>
        <b-col cols="2">
          <b-button @click="decrementPage" :disabled="loading" pill class="shadow-sm">
            <strong>Prev</strong>
          </b-button>
        </b-col>

        <b-col cols="2">
          <b-button @click="incrementPage" :disabled="loading" pill class="shadow-sm">
            <strong>Next</strong>
          </b-button>
        </b-col>

        <b-col class="title" cols="4">
          <strong>{{ fileName }}</strong>
        </b-col>

        <b-col cols="2">
          <b-button @click="decreasePageSize" :disabled="loading" pill class="shadow-sm">
            <strong>Page Size -</strong>
          </b-button>
        </b-col>

        <b-col cols="2">
          <b-button @click="increasePageSize" :disabled="loading" pill class="shadow-sm">
            <strong>Page Size +</strong>
          </b-button>
        </b-col>

      </b-row>
    </b-container>
    <div class="divide"/>
    <div class="content" v-if="!loading">
      <div v-for="line in slice.lines" :key="line.Position" >
        <div :class="(line.Position % 2 == 0) ? 'even' : 'odd'">
          {{ line.Position + ':' + line.Content.toString() }}
        </div>
      </div>
    </div>
    <div class="loading" v-else>
      <b-spinner/>
    </div>
  </div>
</template>

<script>
import FileService from '@/service/file'

export default {
  name: 'text-viewer',
  props: {
    fileName: String,
    finalPage: Number,
    acr: String
  },
  data () {
    return {
      slice: {},
      fileID: '',
      pageSize: 0,
      pageSizeIncrement: 5,
      startPosition: 0
    }
  },
  computed: {
    endPosition () {
      return this.startPosition + this.pageSize
    },
    loading () {
      return this.$store.state.file.isLoading
    }
  },
  methods: {
    fetchSlice () {
      FileService.slice({
        fid: this.fileID,
        start: this.startPosition,
        end: this.endPosition
      }).then(({ data }) => {
        this.slice = data
      })
    },
    decreasePageSize () {
      this.pageSize = Math.max(1, this.pageSize - this.pageSizeIncrement)
      this.fetchSlice()
    },
    increasePageSize () {
      this.pageSize += this.pageSizeIncrement
      if (this.endPosition > this.finalPage) {
        this.startPosition = this.finalPage - this.pageSize
        if (this.startPosition < 0) {
          this.startPosition = 0
          this.pageSize = this.finalPage
        }
      }
      this.fetchSlice()
    },
    incrementPage () {
      this.startPosition += this.pageSize
      if (this.endPosition > this.finalPage) {
        this.startPosition = this.finalPage - this.pageSize
      }
      this.fetchSlice()
    },
    decrementPage () {
      this.startPosition = Math.max(0, this.startPosition - this.pageSize)
      this.fetchSlice()
    }
  },
  created () {
    this.fileID = this.$route.params.id
    this.pageSize = Math.min(10, this.finalPage)
    this.fetchSlice()
  },
  watch: {
    finalPage () {
      this.pageSize = Math.min(10, this.finalPage)
      this.fetchSlice()
    }
  }
}
</script>

<style scoped lang="scss">

.header {
  margin-bottom: 5px;

  .title {
    text-align: center;
    margin-top: auto;
    margin-bottom: auto;
  }
}

.divide {
  margin-top: 2px;
  margin-bottom: 20px;
  height: 1px;
  width: 100%;
  border-top: 1px solid gray;
}

.content {
  margin-left: 20px;
  margin-right: 20px;

  .even {
    background-color: $app-bg1;
  }

  .odd {
    background-color: $app-bg;
  }
}

button {
  width: 100%;
  color: rgb(46, 46, 46);
  background-color: white;
}

.loading {
  margin-top: auto;
  margin-bottom: auto;
  text-align: center;
}

</style>
