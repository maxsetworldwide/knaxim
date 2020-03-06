<template>
  <b-table
    ref="table"
    striped
    hover
    selectable
    :items="fileRows"
    :fields="columnHeaders"
    :busy="busy"
    :sort-compare="sortCompare"
    @row-selected="onCheck"
  >
  <template v-slot:table-colgroup="scope">
    <col
      v-for="field in scope.fields"
      :key="field.key"
      :class="field.class"
    >
  </template>
  <template v-slot:head(expand)="col">
    <!-- <svg @click.stop="expandAll">
      <use href="../assets/app.svg#expand-tri" class="triangle"/>
    </svg> -->
  </template>
  <template v-slot:table-busy>
    <div class="text-center">
      <b-spinner class="align-middle"></b-spinner>
      <strong>Loading...</strong>
    </div>
  </template>
  <template v-slot:head(select)>
    <b-checkbox v-model="selectAllMode"/>
  </template>
  <template v-slot:cell(select)="{ rowSelected }">
    <template v-if="rowSelected">
      <span aria-hidden="true">&check;</span>
    </template>
    <template v-else>
      <span aria-hidden="true">&nbsp;</span>
    </template>
  </template>
  <template v-slot:head(action)>
    <slot name="action"></slot>
  </template>
  <template v-slot:cell(name)="data">
    <!-- <span v-if="data.item.isFolder" class="file-name" @click.prevent.stop="openFolder(data.value)">{{ data.value }}</span> -->
    <span class="file-name" @click="open(data.item.id)">{{ data.value }}</span>
  </template>
  <template v-slot:cell(expand)="row">
    <div @click.stop="openPreview(row)">
      <b-icon v-if="!row.detailsShowing" icon="chevron-down" class="expand"/>
      <b-icon v-else icon="chevron-up" class="expand"/>
    </div>
  </template>
  <template v-slot:row-details="row">
    <b-spinner v-if="filePreview[row.item.id].loading" class="align-middle"></b-spinner>
    <span v-else>{{ filePreview[row.item.id].lines ? filePreview[row.item.id].lines.join(' ') : '' }}</span>
  </template>
  <template v-slot:cell(action)="data">
    <file-icon :extention="(data.item.ext || '')" :webpage="!!data.item.url"/>
  </template>
  </b-table>
</template>
<script>
import fileIcon from '@/components/file-icon'
import { mapGetters, mapActions } from 'vuex'
import { LOAD_OWNER, LOAD_PREVIEW } from '@/store/actions.type'
import { humanReadableSize, humanReadableTime } from '@/plugins/utils'

export default {
  name: 'file-table',
  components: {
    fileIcon
  },
  props: {
    files: {
      type: Array,
      default: () => []
    },
    /*
     * folders: {
     *   type: Array,
     *   default: () => []
     * },
     */
    busy: Boolean
  },
  data () {
    return {
      selected: false,
      columnHeaders: [
        {
          key: 'select'
        },
        {
          key: 'action',
          class: 'action-column'
        },
        {
          key: 'name',
          class: 'name-column',
          sortable: true
        },
        {
          key: 'expand',
          label: '',
          class: 'expand-column'
        },
        {
          key: 'owner',
          sortable: true
        },
        {
          key: 'date',
          sortable: true
        },
        {
          key: 'size',
          sortable: true
        }
      ]
    }
  },
  computed: {
    selectAllMode: {
      get () {
        return this.selected
      },
      set (newValue) {
        if (!newValue && this.selected) {
          this.unselectAll()
        }
        if (newValue && !this.selected) {
          this.selectAll()
        }
      }
    },
    /*
     * folderRows () {
     *   let id = 0
     *   return this.folders.map(name => {
     *     id++
     *     return {
     *       isFolder: true,
     *       name,
     *       id
     *     }
     *   })
     * },
     */
    fileRows () {
      // console.log(this.populateFiles)
      return this.populateFiles(this.files).filter(f => f).map(file => {
        this[LOAD_OWNER]({ id: file.owner })
        let matches = file.name.match(/(?:^(.*)\.([^/.]{1,8})$)|(.*)/)
        let splitname = [matches[1] || matches[0], matches[2] || '']
        return {
          id: file.id,
          url: file.url,
          // isFolder: false,
          name: splitname[0],
          ext: splitname[1],
          owner: this.ownerNames[file.owner],
          size: file.size && humanReadableSize(file.size),
          sizeInt: file.size,
          date: file.date && humanReadableTime(file.date.upload),
          dateInt: file.date ? Date.parse(file.date.upload) : 0,
          _showDetails: (file._showDetails || false)
        }
      })
    },
    anyRowExpanded () {
      return this.fileRows.reduce((acc, row) => {
        return acc || row._showDetails
      }, false)
    },
    ...mapGetters(['ownerNames', 'populateFiles', 'previewLoading', 'filePreview'])
  },
  methods: {
    openPreview (row) {
      row.toggleDetails()
      this[LOAD_PREVIEW](row.item)
    },
    expandALL () {
      const expand = !this.anyRowExpanded
      this.fileRows.forEach((file) => {
        file._showDetails = expand
      })
    },
    unselectAll () {
      this.$refs.table.clearSelected()
    },
    selectAll () {
      this.$refs.table.selectAllRows()
    },
    onCheck (items) {
      // items = items.filter(r => !r.isFolder).map(r => r.id)
      items = items.map(r => r.id)
      if (items.length > 0) {
        this.selected = true
      } else {
        this.selected = false
      }
      this.$emit('selection', items)
    },
    sortCompare (a, b, key) {
      if (key === 'date') {
        return a.dateInt - b.dateInt
      } else if (key === 'size') {
        return a.sizeInt - b.sizeInt
      } else {
        return null
      }
    },
    /*
     * openFolder (name) {
     *   this.$emit('open-folder', name)
     * },
     */
    open (id) {
      this.$emit('open', id)
    },
    // ...mapGetters(['populateFiles']),
    ...mapActions([LOAD_OWNER, LOAD_PREVIEW])
  }
}
</script>
<style lang="scss" scoped>
  .file-name {
    cursor: pointer;
    color: $app-clr1;

    &:hover {
      text-decoration: underline;
    }
  }
  svg {
    width: 100%;
    height: 40px;
  }

  .action-column {
    width: 5%;
  }

  .name-column {
    min-width: 30%;
  }

  .expand-column {
    width: 8%;
  }

  .expand {
    fill: $app-icon;

    &:hover {
      fill: $app-icon-hl;
    }
  }
</style>
