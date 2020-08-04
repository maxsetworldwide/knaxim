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

const PermissionService = {
  permissions ({ id }) {
    return ApiService.get(`perm/file/${id}`).catch(error => {
      throw buildError('PermissionService.permissions', error)
    })
  },

  share ({ id, targets }) {
    return ApiService.post(`perm/file/${id}`, { 'id': targets }).catch(error => {
      throw buildError('PermissionService.share', error)
    })
  },

  stopShare ({ id, targets }) {
    return ApiService.delete(`perm/file/${id}`, { 'id': targets }).catch(error => {
      throw buildError('PermissionService.stopShare', error)
    })
  },
  makePublic ({ fid }) {
    return ApiService.post(`perm/file/${fid}/public`).catch(error => {
      throw buildError('PermissionService.makePublic', error)
    })
  },

  stopPublic ({ fid }) {
    return ApiService.delete(`perm/file/${fid}/public`).catch(error => {
      throw buildError('PermissionService.stopPublic', error)
    })
  }
}

export default PermissionService
