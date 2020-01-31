// Core Vue resources
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

// 3rd party vendors
import BootstrapVue from 'bootstrap-vue'

// local resources
// import KnaxUtils from './plugins/knax-utils'
import ApiService from './service/api'

import './scss/custom.scss'

Vue.config.productionTip = false
Vue.use(BootstrapVue)
// Vue.use(KnaxUtils)

ApiService.init()

/*
 * To use the event bus:
 * import { EventBus } from '@/main'
 * EventBus.$emit('event-name', payload)
 * EventBus.$on('event-name', func)
 * EventBus.$off('event-name')
 * func (payload) {
 *   ...
 * }
 * https://alligator.io/vuejs/global-event-bus/
 *
 * Please be sure to keep track of events you emit within your components.
 */
export const EventBus = new Vue()

// ApiService._vm: Allow emiting error events from root vue instance.
ApiService._vm = new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
