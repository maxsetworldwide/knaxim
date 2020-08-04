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

import ResponseError from '@/error/ResponseError'

export const buildError = function (prefix, error, suffix = '') {
  if (!process.env.VUE_APP_DEBUG) {
    prefix = ''
    suffix = ''
  } else {
    prefix = prefix + ' '
    suffix = ' ' + suffix
  }
  try {
    return new ResponseError(`${prefix}${error.response.data.message}${suffix}`, error.response)
  } catch {
    return new ResponseError(`${prefix}${error.message}${suffix}`)
  }
}
