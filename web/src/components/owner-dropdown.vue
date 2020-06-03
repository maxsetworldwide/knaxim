<!--
  -  a wrapper for <b-form-datalist>, to be used in conjunction with a b-form-input
 -->
<template>
 <b-form-datalist :id="id" :options="listItems" />
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'owner-dropdown',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  computed: {
    listItems () {
      let result = []
      for (let entry in this.ownerNames) {
        result.push(this.ownerNames[entry])
      }
      return result
    },
    ...mapGetters(['ownerNames'])
  }
}
</script>

<style>

</style>
