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
  <b-modal
    ref="folderModal"
    :id="id"
    centered
    hide-footer
    hide-header
    @hidden="onClose"
    content-class="modal-style"
    >
    <b-form @submit.prevent="create">
      <input type="text" v-model="name" placeholder="Enter Name..."/>
      <b-button class="shadow-sm" type="submit">Create</b-button>
    </b-form>
  </b-modal>
</template>
<script>

export default {
  name: 'folder-modal',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  data () {
    return {
      name: ''
    }
  },
  methods: {
    create () {
      this.$emit('new-folder', this.name)
      this.hide()
    },
    onClose () {
      this.$emit('close')
    },
    show () {
      this.$refs['folderModal'].show()
    },
    hide () {
      this.$refs['folderModal'].hide()
    }
  }
}
</script>
