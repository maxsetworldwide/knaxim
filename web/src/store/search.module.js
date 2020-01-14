import Vue from 'vue'
import FileService from '@/service/file'
import SearchService from '@/service/search'
import { FILES_SEARCH, LOAD_FILE_MATCHES } from './actions.type'
import {
  FILES_SEARCH_START,
  UPDATE_SEARCH_HISTORY,
  GET_SLICES,
  ADD_FILE_META,
  UPDATE_FILE_META,
  FILES_SEARCH_END
} from './mutations.type'

const state = {
  files: [],
  isLoading: false,
  slicesLoading: false,
  selected: '',
  history: [],
  maxSummary: 100,
  cancelSearch: false
}

const actions = {
  // TODO: isLoading is updated correctly but not Utalized yet.  If a user spams
  // the FILES_SEARCH action they might get duplicate results.
  // TODO: Use isLoading: true to cancel the current search and start a new one.
  async [FILES_SEARCH] (state, params) {
    if (params.acr) {
      params.find = `"${params.find}" ${params.acr}`
    }
    if (params.find === state.state.selected) {
      return
    }
    state.commit(UPDATE_SEARCH_HISTORY, params)
    state.commit(FILES_SEARCH_START)

    let fileList = await new Promise((resolve, reject) => {
      state.commit('cancelSearch', () => {
        reject(new Error('search canceled'))
      })
      SearchService.userFiles(params).then(({ data }) => {
        if (data.matched && data.matched.length > 0) {
          return data.matched.map(item => {
            return {
              ...item.file,
              count: item.count
            }
          })
        }
        return []
      }).then(r => resolve(r)).catch(e => reject(e))
    })
    state.commit(GET_SLICES)
    fileList.forEach((file) => {
      FileService.search({
        fid: file.id,
        start: 0,
        end: this.state.search.maxSummary,
        find: params.find
      }).then(({ data }) => {
        file.lines = data.lines
        state.commit(ADD_FILE_META, { file })
        state.dispatch(LOAD_FILE_MATCHES, { find: params.find, id: file.id })
      })
    })
    state.commit(FILES_SEARCH_END)
  },
  async [LOAD_FILE_MATCHES] (context, { find, id }) {
    let file = context.getters.searchMatches.reduce((acc, f) => {
      if (f.id === id) {
        return f
      }
      return acc
    }, {})
    for (let start = context.state.maxSummary; start < file.count && (file.lines || []).length < 4; start += context.state.maxSummary) {
      let lines = await FileService.search({
        fid: id,
        start,
        end: start + context.state.maxSummary,
        find
      }).then(({ data }) => data.lines)
      if ((lines || []).length) {
        context.commit(UPDATE_FILE_META, { id, lines })
      }
    }
  }
}

const mutations = {
  [FILES_SEARCH_START] (state) {
    state.files = []
    state.isLoading = true
  },
  [UPDATE_SEARCH_HISTORY] (state, { find }) {
    state.selected = find

    if ((typeof find !== 'string' || find.length < 1) ||
        state.history.find((i) => { return i === find })) {
      return false
    }

    if (state.history.unshift(find) > 10) {
      state.history.pop()
    }
  },
  [GET_SLICES] (state) {
    state.slicesloading = true
  },
  [ADD_FILE_META] (state, { file }) {
    state.files.push(file)
  },
  [UPDATE_FILE_META] (state, { id, lines }) {
    state.files.forEach((f, i, arr) => {
      if (f.id === id) {
        Vue.set(arr[i], 'lines', (f.lines || []).concat(lines))
      }
    })
  },
  [FILES_SEARCH_END] (state) {
    state.slicesLoading = false
    state.isLoading = false
  },
  cancelSearch (state, newfunc) {
    if (state.cancelSearch) {
      state.cancelSearch()
    }
    state.cancelSearch = newfunc
  }
}

const getters = {
  searchMatches (state) {
    return state.files
  },
  currentSearch (state) {
    return state.selected
  },
  searchHistory (state) {
    return state.history
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
