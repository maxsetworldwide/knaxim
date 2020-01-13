import FileService from '@/service/file'
import { FILE_SLICES, CREATE_FILE } from './actions.type'
import {
  START_SLICES,
  END_SLICES,
  FILE_CREATED
} from './mutations.type'

const state = {
  fileId: '',
  lines: [],
  size: 0,
  isLoading: false
}

const actions = {
  [CREATE_FILE] (context, params) {
    return new Promise(resolve => {
      FileService.create(params).then(({ data }) => {
        resolve(data)
        context.commit(FILE_CREATED, data.id)
      }).catch((error) => {
        throw new Error(error)
      })
    })
  },

  /**
   * Get a set of sentences from a file.
   *
   * @param {object} context  state
   * @param {object} params  url properties
   * @param {string} params.fid  file_id
   * @param {number} params.start  starting sentence index
   * @param {number} params.end  ending sentence index
   * @return {Promise}
   */
  [FILE_SLICES] (context, params) {
    context.commit(START_SLICES)
    return new Promise((resolve, reject) => {
      FileService.slice(params)
        .then(({ data }) => {
          context.commit(END_SLICES, {
            slices: data.lines,
            count: data.size
          })
          resolve(data)
        }).catch((error) => {
          context.commit(END_SLICES)
          reject(error)
        })
    })
  }
}

const mutations = {
  [FILE_CREATED] (state, fileId) {
    state.fileId = fileId
  },

  [START_SLICES] (state) {
    state.isLoading = true
  },
  [END_SLICES] (state, { slices, count }) {
    state.lines = slices
    state.size = count
    state.isLoading = false
  }
}

const getters = {
}

export default {
  state,
  actions,
  mutations,
  getters
}
