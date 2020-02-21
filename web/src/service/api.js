import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

const ApiService = {
  _vm: null,

  /**
   * init - Set default http client settings.
   */
  init () {
    Vue.use(VueAxios, axios)
    // config values are set in .env.[dev, test, production] files.
    // they can be overridden with a similar file with a .local extension.
    Vue.axios.defaults.baseURL = process.env.VUE_APP_API_URL
    Vue.axios.defaults.withCredentials = true
    Vue.axios.defaults.transformRequest = [
      function (data, header) {
        let fdata = new FormData()
        for (let key in data) {
          if (typeof data[key] !== 'undefined') {
            fdata.append(key, data[key])
          }
        }
        return fdata
      }
    ]

    Vue.axios.defaults.transformResponse = [
      function (data) {
        try {
          return JSON.parse(data)
        } catch {
          return data
        }
      }
    ]
  },

  query (resource, params) {
    return Vue.axios.get(resource, { params: params })
  },

  get (resource, slug = '') {
    return Vue.axios.get(`${resource}` + (slug ? `/${slug}` : ''))
  },

  post (resource, params) {
    return Vue.axios.post(`${resource}`, params)
  },

  put (resource, params) {
    return Vue.axios.put(`${resource}`, params)
  },

  delete (resource, params) {
    return Vue.axios.delete(resource, { data: params })
  }
}

export default ApiService
