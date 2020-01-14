<template>
  <svg v-if="hasPic">
    <use v-if="folder" href="../assets/app.svg#folder"></use>
    <use v-else-if="webpage" href="../assets/app.svg#webpage"></use>
    <use v-else-if="ext === 'pdf'" href="../assets/app.svg#pdf2"></use>
    <use v-else-if="ext === 'doc' || ext === 'docx'" href="../assets/app.svg#doc"></use>
    <use v-else-if="ext === 'csv'" href="../assets/app.svg#csv"></use>
    <use v-else-if="ext === 'txt'" href="../assets/app.svg#txt"></use>
    <use v-else-if="ext === 'ppt' || ext === 'pptx'" href="../assets/app.svg#ppt"></use>
    <use v-else-if="ext === 'xls' || ext === 'xlsx'" href="../assets/app.svg#xls"></use>
  </svg>
  <h5 class="unrecognized" v-else>{{ extUpper }}</h5>
</template>

<script>

export default {
  name: 'file-icon',
  props: {
    extention: String,
    folder: Boolean,
    webpage: Boolean
  },
  computed: {
    ext () {
      return this.extention.toLowerCase()
    },
    hasPic () {
      return ['pdf', 'doc', 'docx', 'csv', 'ppt', 'pptx', 'txt', 'xls', 'xlsx'].reduce((acc, e) => acc || e === this.ext, this.folder || this.webpage)
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
