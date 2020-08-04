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
import { shallowMount, createLocalVue } from '@vue/test-utils'
import FileList from '@/components/file-list'
import Vuex from 'vuex'
import {
  LOAD_FOLDERS,
  PUT_FILE_FOLDER,
  REMOVE_FILE_FOLDER,
  GET_USER
} from '@/store/actions.type'

const localVue = createLocalVue()
localVue.use(Vuex)

let actions = {
  [LOAD_FOLDERS] () {
    return true
  },
  [PUT_FILE_FOLDER] () {
    return true
  },
  [REMOVE_FILE_FOLDER] () {
    return true
  },
  [GET_USER] () {
    return true
  }
}

let store = new Vuex.Store({
  actions,
  getters: {
    activeFolders () {
      return []
    }
  }
})

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(FileList, {
    stubs: [
      'b-table',
      'b-tooltip',
      'upload-modal',
      'file-actions',
      'folder-modal',
      'share-modal',
      'file-icon',
      'file-table'
    ],
    store,
    localVue,
    watch: {
      gid () {
        return ''
      }
    },
    propsData: {
      ...options.props
    },
    methods: {
      showAuth () {
        return true
      },
      refresh () {
        return true
      },
      adjustFavorite (add) {
        add += ''
        return true
      },
      createFolder (name) {
        name += ''
        return true
      },
      open (id) {
        id += ''
        return true
      },
      ...options.methods
    },
    computed: {
      fileMap () {
        return true
      },
      activeGroup () {
        return true
      },
      promptUpload () {
        return false
      },
      isAuthenticated () {
        return true
      },
      gid () {
        return ''
      },
      files () {
        return []
      },
      filterFiles () {
        return true
      },
      fileRows () {
        return []
      },
      folders () {
        return true
      },
      favoriteFolder () {
        return true
      },
      selected () {
        return true
      },
      selectedAllMode () {
        return true
      },
      anyRowExpanded () {
        return {}
      },
      ...options.computed
    }
  })
}

describe('FileList', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(FileList)).toBe(true)
  })
})
