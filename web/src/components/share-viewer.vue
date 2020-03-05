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
