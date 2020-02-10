// Core Vue resources
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

// 3rd party vendors
import BootstrapVue from 'bootstrap-vue'

// local resources
import ApiService from './service/api'

import './scss/custom.scss'

Vue.config.productionTip = false
Vue.use(BootstrapVue)

ApiService.init()

// ApiService._vm: Allow emiting error events from root vue instance.
ApiService._vm = new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
