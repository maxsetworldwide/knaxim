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
// import ApiService, { FileService, FilesService } from '@/common/api.service'

import ApiService from '@/service/api'
import UserService from '@/service/user'
import FileService from '@/service/file'
import SearchService from '@/service/search'

export async function getDigestLast8 (message) {
  const encoder = new TextEncoder()
  const arrayBuffer = await crypto.subtle.digest('SHA-1', encoder.encode(message))
  const array = Array.from(new Uint8Array(arrayBuffer))
  return array.map(item => item.toString(16).padStart(2, '0'))
    .join('').substr(-8)
}

export default async function setDefaults (that) {
  // Create some default test values.
  Object.assign(that, {
    login: '[TEST_SETUP]',
    email: '[TEST_SETUP]',
    password: 'testErr1',

    groupId: 'aTestGroup',
    groupName: '',

    file: new Blob(['<a id="a"><p>I have seen...all good people...and I like it.</p><b id="b">And then, and that!</b></a>'],
      { type: 'text/html' }),
    fileId: '[TEST_SETUP]',
    fileName: 'blob',
    find: 'and',

    directoryId: 'KKwUc_hx',
    directoryName: 'aTestDirectory'
  })

  ApiService.init()

  let digest = await getDigestLast8(window.navigator.userAgent)
  Object.assign(that, {
    login: digest,
    email: `${digest}@example.com`
  })

  // Create a user, ignore duplicate-user error.
  await UserService.create({
    email: that.email,
    name: that.login,
    password: that.password
  }).catch((error) => {
    error.message.indexOf('409') > 0 ||
    fail(`Test Suite Setup: ${error.message}`)
  })

  // Login
  await UserService.login({
    name: that.login,
    pass: that.password
  }).catch((error) => {
    fail(`Test Suite Setup: userlogin ${error.message}`)
  })

  // Create some files.
  //  - OR -
  // Store a file ID from an existing file.
  await SearchService.userFiles({ find: that.find }).then(({ data }) => {
    that.fileId = data.matched[0].file.id
  }).catch(({ message }) => {
    FileService.create({ file: that.file }).then(({ data }) => {
      that.fileId = data.id
      FileService.create({ file: that.file }).catch(({ message }) => {
        fail(message)
      })
    }).catch(({ message }) => {
      fail(message)
    })
  })

  console.log(`Email: ${that.email}   Login: ${that.login}   Password: ${that.password}`)
}
