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
import { TOUCH } from '@/store/mutations.type'
import modul from '@/store/recents.module'

describe('Recent Module', function () {
  beforeAll(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  beforeEach(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  describe('Mutations', function () {
    const m = modul.mutations
    it('Records files', function () {
      m[TOUCH](this.state, 'fid')
      expect(this.state.files).toEqual(['fid'])
    })
    it('touched files are moved to the front', function () {
      this.state.files = ['1', '2']
      m[TOUCH](this.state, '2')
      expect(this.state.files).toEqual(['2', '1'])
    })
  })
})
