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

const GroupService = {
  create ({ name }) {
    return ApiService.put(`group`, { 'newname': name }).catch(error => {
      throw buildError('GroupService.create', error)
    })
  },

  // get associated groups, provide gid to get groups gid belongs to
  associated ({ gid }) {
    return ApiService.get(`group/options${gid ? '/' + gid : ''}`).catch(error => {
      throw buildError('GroupService.associated', error)
    })
  },

  info ({ gid }) {
    return ApiService.get(`group/${gid}`).catch(error => {
      throw buildError('GroupService.info', error)
    })
  },

  add ({ gid, target }) {
    return ApiService.post(`group/${gid}/member`, { 'id': target }).catch(error => {
      throw buildError('GroupService.add', error)
    })
  },

  remove ({ gid, target }) {
    return ApiService.delete(`group/${gid}/member`, { 'id': target }).catch(error => {
      throw buildError('GroupService.remove', error)
    })
  },

  lookup ({ name }) {
    return ApiService.get(`group/name/${name}`).catch(error => {
      throw buildError('GroupService.lookup', error)
    })
  }
}

export default GroupService
