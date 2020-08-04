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

const OwnerService = {
  id (id) {
    return ApiService.get(`owner/id/${id}`).catch(error => {
      throw buildError('OwnerService.id', error)
    })
  },
  name (name) {
    return ApiService.get(`owner/name/${name}`).catch(error => {
      throw buildError('OwnerService.name', error)
    })
  }
}

export default OwnerService
