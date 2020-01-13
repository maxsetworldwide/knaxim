import FileService from '@/service/file'
import { FILES_LIST } from './actions.type'
import {
  FETCH_FILES_START,
  FETCH_FILES_END
} from './mutations.type'

const state = {
  files: {},
  isLoading: false
}

const actions = {
  async [FILES_LIST] (context, params) {
    context.commit(FETCH_FILES_START)

    // TODO: This API should return an object...It is working with an array ATM.
    var res = await FileService.list(params || {})
    context.commit(FETCH_FILES_END, {
      files: Object.values(res.data.files).reduce((obj, fileinfo) => {
        obj[fileinfo.file.id] = {
          ...fileinfo.file,
          count: fileinfo.count,
          size: fileinfo.size
        }
        return obj
      }, {})
    })
  }
}

const mutations = {
  [FETCH_FILES_START] (state) {
    state.isLoading = true
  },
  [FETCH_FILES_END] (state, { folders, files }) {
    state.files = { ...state.files, ...files }
    state.isLoading = false
  }
}

const getters = {
  fileMap (state) {
    return state.files
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
