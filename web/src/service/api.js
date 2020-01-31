import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

const ApiService = {
  // TODO: Is there another/better way to access the root Vue instance?!
  // Expose the root vue instance for...
  //  - sending global error events from ApiService (app::set::error)
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

    // *** Interceptors appear to be added in a stack, LIFO.
    // *** They also fail silently!

    // response interceptor to attempt conversion of all json data.
    // TODO: Remove this interceptor when the server starts sending the
    // correct content-type.
    Vue.axios.interceptors.response.use(function (response) {
      try {
        let data = JSON.parse(response.data)
        response.data = data
      } catch (e) {}
      return response
    }, function (error) {
      return Promise.reject(error)
    })

    // Handle global API response error messages.
    Vue.axios.interceptors.response.use(function (config) {
      if (config.data && config.data.message === 'login') {
        ApiService._vm.$root.$emit('app::set::error', 'Please Login.')
        // TODO: Throwing an error here blindly prevents duplicate requests
        // in some cases.  This may need to be handled better by the code
        // making requests, or perhaps a different mechanisim is needed.
        throw new Error('Login Required')
      }
      return config
    })

    // TODO: Remove this code when we are confident that objects are accepted
    //  by the server;  It will clean things up a bit.  This might be a
    //  performance decision.

    // Request FromData: Convert object into FormData object,
    //  remove undefined params.
    Vue.axios.interceptors.request.use(function (config) {
      if (config.data) {
        let data = new FormData()
        for (let key in config.data) {
          if (typeof config.data[key] !== 'undefined') {
            data.append(key, config.data[key])
          }
        }
        config.data = data
      }

      return config
    })
  },

  query (resource, params) {
    return Vue.axios.get(resource, { params: params }).catch(error => {
      throw new Error(`ApiService ${error}`)
    })
  },

  get (resource, slug = '') {
    return Vue.axios.get(`${resource}` + (slug ? `/${slug}` : ''))
      .catch(error => {
        throw new Error(`ApiService ${error}`)
      })
  },

  post (resource, params) {
    return Vue.axios.post(`${resource}`, params)
  },

  put (resource, params) {
    return Vue.axios.put(`${resource}`, params)
  },

  delete (resource, params) {
    return Vue.axios.delete(resource, { data: params }).catch(error => {
      throw new Error(`ApiService ${error}`)
    })
  }
}

export default ApiService
