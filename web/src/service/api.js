import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

const ApiService = {
  _vm: null,

  /**
   * init - Set default http client settings.
   */
  init () {
    // config values are set in .env.[dev, test, production] files.
    // they can be overridden with a similar file with a .local extension.
    axios.defaults.baseURL = process.env.VUE_APP_API_URL
    axios.defaults.withCredentials = true
    axios.defaults.transformRequest = [
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

    axios.defaults.transformResponse = [
      function (data) {
        try {
          return JSON.parse(data)
        } catch {
          return data
        }
      }
    ]
    Vue.use(VueAxios, axios)
  },

  query (resource, params) {
    return axios.get(resource, { params: params })
  },

  get (resource, slug = '') {
    return axios.get(`${resource}` + (slug ? `/${slug}` : ''))
  },

  post (resource, params) {
    return axios.post(`${resource}`, params)
  },

  put (resource, params) {
    return axios.put(`${resource}`, params)
  },

  delete (resource, params) {
    return axios.delete(`${resource}`, { data: params })
  }
}

export default ApiService
