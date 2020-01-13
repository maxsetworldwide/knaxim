import { TOUCH } from './mutations.type'

const state = {
  files: []
}

const mutations = {
  [TOUCH] (state, fileID) {
    var newfiles = state.files.filter((val) => { return val !== fileID })
    newfiles.unshift(fileID)
    state.files = newfiles
  }
}

export default {
  state,
  mutations
}
