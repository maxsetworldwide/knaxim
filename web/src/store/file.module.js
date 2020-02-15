import FileService from '@/service/file'
import { FILE_SLICES, CREATE_FILE } from './actions.type'
import {
  START_SLICES,
  END_SLICES,
  FILE_CREATED,
  FILE_START_LOADING,
  FILE_STOP_LOADING
} from './mutations.type'

const state = {
  loading: 0
}

const actions = {
  [CREATE_FILE] (context, params) {
    context.commit(FILE_START_LOADING)
    return FileService.create(params)
      .then(res => res.data)
      .finally(() => { context.commit(FILE_STOP_LOADING) })
  }
}

const mutations = {
  [FILE_START_LOADING] (state) {
    state.loading += 1
  },
  [FILE_STOP_LOADING] (state) {
    state.loading -= 1
  }
}

const getters = {
  fileLoading (state) {
    return state.loading > 0
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
