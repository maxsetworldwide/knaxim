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
confirm-modal: used for relaying feedback to the user and receiving input

This is a modal that is called via a global function call. Its data is
initialized by the function call (there are no props, only data), and returns
a promise based on the user's response. The button with okText will resolve
to true, and the button with cancelText as well as exiting the modal will
resolve to false.

Any of the data fields can be changed via the object you pass to the utility
function call.

See knax-utils for more info
-->
<template>
  <b-modal
  :id="id"
  ref="modal"
  @hidden="onClose"
  :no-close-on-backdrop="noCloseOnBackdrop"
  :no-close-on-esc="noCloseOnEsc"
  centered
  hide-footer
  hide-header
  content-class="modal-style">
    {{ msg }}
    <div v-if="okText.length > 0">
      <b-button class="shadow-sm" @click="onYes">
        {{ okText }}
      </b-button>
    </div>
    <div v-if="cancelText.length > 0">
      <b-button class="shadow-sm" @click="onCancel">
        {{ cancelText }}
      </b-button>
    </div>
  </b-modal>
</template>

<script>

export default {
  name: 'confirm-modal',
  data () {
    return {
      id: '',
      promResolve: null,
      msg: '',
      okText: '',
      cancelText: '',
      noCloseOnBackdrop: false,
      noCloseOnEsc: false
    }
  },
  methods: {
    onClose () {
      if (this.promResolve) {
        this.promResolve(false)
      }
      this.promResolve = null
    },
    onYes () {
      this.promResolve(true)
      this.$refs['modal'].hide()
    },
    onCancel () {
      this.promResolve(false)
      this.$refs['modal'].hide()
    },
    promise () {
      return new Promise((resolve) => {
        this.promResolve = resolve
      })
    },
    show () {
      this.$refs['modal'].show()
    }
  },
  mounted () {
    this.show()
  }
}
</script>

<style scoped lang="scss">

button {
  @extend %pill-buttons;
  width: 20%;
  margin-right: 5px;
  margin-left: 5px;
  margin-top: 5px;
}

::v-deep .modal-style {
  @extend %modal-corners;
  text-align: center;
}

</style>
