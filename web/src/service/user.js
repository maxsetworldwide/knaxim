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

import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const UserService = {
  create ({ name, password, email }) {
    return ApiService.put(`user`, { 'name': name, 'pass': password, 'email': email }).catch(error => {
      throw buildError('UserService.create', error)
    })
  },

  createAdmin ({ name, password, email, adminKey }) {
    return ApiService.put(`user/admin`, { 'name': name, 'pass': password, 'email': email, 'adminKey': adminKey }).catch(error => {
      throw buildError('UserService.createAdmin', error)
    })
  },

  info ({ id }) {
    if (id === undefined) {
      return ApiService.get(`user`).catch(error => {
        throw buildError('UserService.info', error)
      })
    } else {
      return ApiService.query(`user`, { 'id': id }).catch(error => {
        throw buildError(`UserService.info{${id}}`, error)
      })
    }
  },

  lookup ({ name }) {
    return ApiService.get(`user/name/${name}`).catch(error => {
      throw buildError('UserService.lookup', error)
    })
  },

  login ({ name, pass }) {
    return ApiService.post(`user/login`, { 'name': name, 'pass': pass }).catch(error => {
      throw buildError('UserService.login', error)
    })
  },

  logout () {
    return ApiService.delete(`user`).catch(error => {
      throw buildError('UserService.logout', error)
    })
  },

  completeProfile () {
    return ApiService.get(`user/complete`).catch(error => {
      throw buildError('UserService.completeProfile', error)
    })
  },

  changePassword ({ oldpass, newpass }) {
    return ApiService.post(`user/pass`, { 'oldpass': oldpass, 'newpass': newpass }).catch(error => {
      throw buildError('UserService.changePassword', error)
    })
  },

  data () {
    return ApiService.get(`user/data`).catch(error => {
      throw buildError('UserService.data', error)
    })
  },

  requestReset ({ name }) {
    return ApiService.put('user/reset', { name }).catch(error => {
      throw buildError('UserService.requestReset', error)
    })
  },

  resetPass ({ passkey, newpass }) {
    return ApiService.post('user/reset', { 'key': passkey, 'newpass': newpass }).catch(error => {
      throw buildError('UserService.resetPass', error)
    })
  }
}

export default UserService
