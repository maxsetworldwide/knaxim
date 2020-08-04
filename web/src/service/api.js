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
        if (data && data._sendJSON) {
          return JSON.stringify(data)
        } else {
          let fdata = new FormData()
          for (let key in data) {
            if (typeof data[key] !== 'undefined') {
              fdata.append(key, data[key])
            }
          }
          return fdata
        }
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
