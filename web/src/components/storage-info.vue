<template>
  <div class="storage-info" v-if="!!this.currentUser.data">
    <div class="pic">
      <svg>
        <use href="../assets/app.svg#server"/>
      </svg>
    </div>
    <div class="data">
      <p class="barlabel">My Storage</p>
      <b-progress :max="total">
        <b-progress-bar :value="current" :variant="variant">
          <span class="label">{{ label }}</span>
        </b-progress-bar>
      </b-progress>
    </div>
  </div>
</template>

<script>
import { GET_USER } from '@/store/actions.type'
import { mapGetters } from 'vuex'
import { humanReadableSize } from '@/plugins/utils'
import { EventBus } from '@/main'

export default {
  name: 'storage-info',
  props: {
  },
  created () {
    EventBus.$on(['file-upload', 'url-upload'], this.refresh)
  },
  method: {
    refresh () {
      this.$store.dispatch(GET_USER)
    }
  },
  computed: {
    total () {
      if (!this.currentUser.data) {
        return 0
      }
      return this.currentUser.data.total
    },
    totalStr () {
      return humanReadableSize(this.total)
    },
    current () {
      if (!this.currentUser.data) {
        return 0
      }
      return this.currentUser.data.current
    },
    currentStr () {
      return humanReadableSize(this.current)
    },
    label () {
      return `${this.currentStr} / ${this.totalStr}`
    },
    variant () {
      if (this.current > this.total * 0.85) {
        return 'danger'
      }
      if (this.current > this.total * 0.65) {
        return 'warning'
      }
      return 'success'
    },
    ...mapGetters(['currentUser'])
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.storage-info {
  max-width: 100%;
  display: grid;
  grid-template-columns: 35px auto;
  grid-template-rows: auto;
  grid-template-areas:
    "pic data";
}

.pic {
  grid-area: pic;
  margin: 2px;
}

.data {
  grid-area: data;
  margin: 2px;
}

.label {
  color: black;
  margin-left: 1em;
}

.barlabel {
  margin: 0;
  font-size: 12px;
}

svg {
  width: 100%;
  height: 35px;

}

</style>
