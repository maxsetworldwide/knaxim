<template>
  <b-dropdown class="file-list-batch" no-caret variant="link" size="sm">
    <template v-slot:button-content>
      <svg class="more">
        <use href="../assets/app.svg#more"/>
      </svg>
    </template>

    <b-dropdown-item href="#" :disabled="!fileSelected" @click="newFolder">
      <svg>
        <use href="../assets/app.svg#folder-2"/>
      </svg>
      <span>Folder+</span>
    </b-dropdown-item>

    <b-dropdown-divider/>

    <b-dropdown-item href="#" :disabled="!fileSelected" @click="share">
      <svg>
        <use href="../assets/app.svg#share"/>
      </svg>
      <span>Share</span>
    </b-dropdown-item>

    <b-dropdown-item href="#" :disabled="!fileSelected" @click="addFavorite">
      <svg>
        <use href="../assets/app.svg#star"/>
      </svg>
      <span v-if="removeFavorite">UnFavorite</span>
      <span v-else>Favorite</span>
    </b-dropdown-item>

    <batch-delete v-if="!singleFile" :files="checkedFiles" :permanent="permanentDelete" #default="{ inputEvents }"
        v-on:delete-files="$emit('delete-files')">
      <b-dropdown-item href="#" v-on="inputEvents" :disabled="!fileSelected">
        <svg>
          <use href="../assets/app.svg#bin"/>
        </svg>
        <span>Delete</span>
      </b-dropdown-item>
    </batch-delete>

    <!-- <b-dropdown-divider/>
    <b-dropdown-item href="#" v-bind:disabled="fileSelected">
      <svg>
        <use href="../assets/app.svg#files"/>
      </svg>
      <span>Redact</span>
    </b-dropdown-item> -->
  </b-dropdown>
</template>

<script>
import BatchDelete from '@/components/batch-delete'

export default {
  name: 'file-list-batch',
  components: {
    BatchDelete
  },
  props: {
    fileSelected: Boolean,
    removeFavorite: Boolean,
    checkedFiles: Array,
    singleFile: Boolean
  },
  methods: {
    newFolder () {
      this.$emit('add-folder')
    },
    addFavorite () {
      this.$emit('favorite', !this.removeFavorite)
    },
    share () {
      this.$emit('share-file')
    }
  },
  computed: {
    trashFolder () {
      return this.$store.state.folder.user['_trash_'] || []
    },
    permanentDelete () {
      return this.checkedFiles.reduce((acc, file) => {
        return acc && this.trashFolder.reduce((a, id) => {
          return a || id === file.id
        }, false)
      }, true)
    }
  }
}
</script>

<style lang="scss">
  .file-list-batch {
    .dropdown {
      height: 35px;
    }

    svg {
      width: 25px;
      height: 25px;
      margin-right: 15px;
    }

    .more {
      width: 40px;
      height: 40px;
      margin-left: 50%;
      margin-top: -20%;
    }
  }
</style>
