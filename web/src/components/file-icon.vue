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
  <svg v-if="hasPic">
    <!-- <use v-if="folder" href="../assets/app.svg#folder"></use> -->
    <use v-if="webpage" href="../assets/app.svg#webpage"></use>
    <use v-else-if="ext === 'pdf'" href="../assets/app.svg#pdf2"></use>
    <use
      v-else-if="ext === 'doc' || ext === 'docx'"
      href="../assets/app.svg#doc"
    ></use>
    <use v-else-if="ext === 'csv'" href="../assets/app.svg#csv"></use>
    <use v-else-if="ext === 'txt'" href="../assets/app.svg#txt"></use>
    <use
      v-else-if="ext === 'ppt' || ext === 'pptx'"
      href="../assets/app.svg#ppt"
    ></use>
    <use
      v-else-if="ext === 'xls' || ext === 'xlsx'"
      href="../assets/app.svg#xls"
    ></use>
  </svg>
  <h5 class="unrecognized" v-else>{{ extUpper }}</h5>
</template>

<script>
export default {
  name: 'file-icon',
  props: {
    extention: String,
    // folder: Boolean,
    webpage: Boolean
  },
  computed: {
    ext () {
      return this.extention.toLowerCase()
    },
    hasPic () {
      return [
        'pdf',
        'doc',
        'docx',
        'csv',
        'ppt',
        'pptx',
        'txt',
        'xls',
        'xlsx'
      ].reduce((acc, e) => acc || e === this.ext, this.webpage)
    },
    extUpper () {
      return this.ext.toUpperCase()
    }
  }
}
</script>

<style scoped>
svg {
  width: 100%;
  height: 40px;
}

.unrecognized {
  text-align: right;
  margin-right: 7px;
}
</style>
