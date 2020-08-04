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
// Core Vue resources
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

// 3rd party vendors
import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'

// local resources
import ApiService from './service/api'

import './scss/custom.scss'

Vue.config.productionTip = false
Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)

ApiService.init()

// ApiService._vm: Allow emiting error events from root vue instance.
ApiService._vm = new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
