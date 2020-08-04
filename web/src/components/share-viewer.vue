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
share-viewer: renders a name with a button to stop sharing

props:
  name: Username of the viewer
  uid: User ID of the viewer
  fids: array of file IDs to operate upon

events:
  'stop-share': emitted upon successful stop share, passing uid
-->
<template>
  <div>
    <span>{{ name }} </span>
    <b-spinner small v-if="loading"/>
    <b-badge v-else
      variant="danger" class="x"
      @click="onClick" href="#"
      v-b-tooltip.hover.right.ds500.dh50 :title="'Stop sharing with ' + name">X</b-badge>
  </div>
</template>

<script>
import PermissionService from '@/service/permission'
import { LOAD_SERVER } from '@/store/actions.type'

export default {
  name: 'share-viewer',
  props: {
    name: String,
    uid: String,
    fids: Array
  },
  data () {
    return {
      loading: false
    }
  },
  methods: {
    onClick () {
      this.loading = true
      Promise.all(this.fids.map(fid => {
        return PermissionService.stopShare({
          id: fid,
          targets: this.uid
        })
      })).then(res => {
        return this.$store.dispatch(LOAD_SERVER)
      }).finally(() => {
        this.loading = false
        this.$emit('stop-share', this.uid)
      })
    }
  }
}
</script>

<style scoped lang="scss">

.x {
  &:hover {
    cursor: pointer;
  }
}

</style>
