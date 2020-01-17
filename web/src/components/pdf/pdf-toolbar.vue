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
  <b-row>
    <b-col offset-md="5" cols="1">
      <b-button @click="increaseScale">
        <svg>
          <use href="@/assets/app.svg#zoom-in"></use>
        </svg>
      </b-button>
    </b-col>
    <b-col cols="1">
      <b-button @click="decreaseScale">
        <svg>
          <use href="@/assets/app.svg#zoom-out"></use>
        </svg>
      </b-button>
    </b-col>
    <b-col cols="2">
      <input
        :value="currPage"
        @input="onPageInput"
        min="1"
        :max="maxPages"
        type="number"
      />
      <span> / {{ maxPages }}</span>
    </b-col>

    <b-col offset-md="2" cols="1">
      <file-list-batch
        fileSelected
        singleFile
        :checkedFiles="mockFileArray"
        :removeFavorite="isFavorite"
        @favorite="adjustFavorite"
        @add-folder="showFolderModal"
        @share-file="showShareModal"
      />
      <folder-modal
        ref="folderModal"
        id="newFolderModal"
        @new-folder="createFolder"
      />
      <share-modal
        hideList
        ref="shareModal"
        id="file-list-share-modal"
        :files="mockFileArray"
      />
    </b-col>
  </b-row>
</template>

<script>
import FileListBatch from '@/components/file-list-batch'
import FolderModal from '@/components/modals/folder-modal'
import ShareModal from '@/components/modals/share-modal'
import { PUT_FILE_FOLDER, REMOVE_FILE_FOLDER } from '@/store/actions.type'
import { mapGetters } from 'vuex'

export default {
  name: 'pdf-toolbar',
  components: {
    FileListBatch,
    ShareModal,
    FolderModal
  },
  props: {
    currPage: Number,
    maxPages: Number,
    id: String
  },
  data () {
    return {
      pageInput: this.currPage
    }
  },
  computed: {
    isFavorite () {
      let favorite = this.$store.state.folder.user['_favorites_'] || []
      return favorite.reduce((acc, id) => {
        return acc || id === this.id
      }, false)
    },
    // TODO: File-List-Batch should be a File-List-Actions (my bad)
    // TODO: Get the actual file information for File-List-Batch.  I'd like to
    //  replace this with a file control component.
    mockFileArray () {
      return [
        {
          id: this.id,
          name: this.id,
          own: '-',
          date: { upload: '-' }
        }
      ]
    },
    ...mapGetters(['activeGroup'])
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
    // file actions
    adjustFavorite (add) {
      return this.$store.dispatch(add ? PUT_FILE_FOLDER : REMOVE_FILE_FOLDER, {
        fid: this.id,
        name: '_favorites_'
      })
    },
    showShareModal () {
      this.$refs['shareModal'].show()
    },
    showFolderModal () {
      this.$refs['folderModal'].show()
    },
    createFolder (name) {
      this.$store.dispatch(PUT_FILE_FOLDER, {
        fid: this.id,
        name,
        group: this.activeGroup ? this.activeGroup.id : undefined
      })
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
</style>
