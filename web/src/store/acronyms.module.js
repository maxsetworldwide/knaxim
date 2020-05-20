import AcronymService from '@/service/acronym'
import { ACRONYMS } from './actions.type'
import {
  SET_ACRONYMS,
  LOADING_ACRONYMS,
  PUSH_ERROR
} from './mutations.type'

const state = {
  acronyms: [],
  loading: 0
}

const getters = {
  acronymResults (state) {
    return state.acronyms
  },
  acronymLoading (state) {
    return state.loading > 0
  }
}

const actions = {
  [ACRONYMS] ({ commit }, { acronym }) {
    if (typeof acronym !== 'string' || acronym.length < 1) {
      commit(SET_ACRONYMS, { acronyms: [] })
      return
    }
    commit(LOADING_ACRONYMS, 1)
    return AcronymService.get({ acronym }).then(data => {
      const { matched } = data.data || []
      commit(SET_ACRONYMS, { acronyms: matched })
      return matched
    })
      .catch(err => commit(PUSH_ERROR, err.addDebug('action ACRONYMS')))
      .finally(() => commit(LOADING_ACRONYMS, -1))
  }
}

const mutations = {
  [SET_ACRONYMS] (state, { acronyms }) {
    state.acronyms = acronyms
  },
  [LOADING_ACRONYMS] (state, delta) {
    state.loading += delta
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
