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
  <div class="h-100">
    <pdf-doc
      v-if="viewExists"
      :fileID="id"
      :acr="acr"
      @no-view="viewExists = false"
    />
    <image-viewer
      v-else-if="imageExists"
      :id="id"
      @no-image="imageExists = false"
    />
    <text-viewer
      v-else
      :fileName="name"
      :finalPage="sentenceCount"
      :acr="acr"
    />
  </div>
</template>

<script>
import TextViewer from '@/components/text-viewer'
import PDFDoc from '@/components/pdf/pdf-doc'
import ImageViewer from '@/components/image-viewer'
import { TOUCH } from '@/store/mutations.type'
import { GET_FILE } from '@/store/actions.type'

export default {
  name: 'file-viewer',
  components: {
    'text-viewer': TextViewer,
    'pdf-doc': PDFDoc,
    'image-viewer': ImageViewer
  },
  props: {
    id: String,
    acr: String
  },
  data () {
    return {
      name: '',
      sentenceCount: 0,
      viewExists: true,
      imageExists: true
    }
  },
  computed: {
    fileType () {
      if (!this.name) {
        return 'txt'
      }
      return this.name
        .split('.')
        .slice(-1)[0]
        .toLowerCase()
    }
  },
  methods: {
    refresh () {
      this.viewExists = true
      this.imageExists = true
      // Add fileID to Recents
      this.$store.commit(TOUCH, this.id)
      this.$store.dispatch(GET_FILE, this).then(data => {
        const { count, name } = data
        this.name = name
        this.sentenceCount = count
      })
    }
  },
  created () {
    this.refresh()
  },
  watch: {
    id (to, from) {
      if (to !== from) {
        this.refresh()
      }
    }
  }
}
</script>

<style></style>
