import AcronymService from '@/service/acronym'
import { ACRONYMS } from './actions.type'
import {
  SET_ACRONYMS,
  BEGIN_ACRONYMS,
  END_ACRONYMS
} from './mutations.type'

const state = {
  acronyms: [],
  loading: false
}

const getters = {
  acronymResults (state) {
    return state.acronyms
  }
}

const actions = {
  [ACRONYMS] (state, { acronym }) {
    if (typeof acronym !== 'string' || acronym.length < 1) {
      state.commit(SET_ACRONYMS, { acronyms: [] })
      return
    }
    return new Promise(resolve => {
      state.commit(BEGIN_ACRONYMS)

      AcronymService.get({ acronym }).then(data => {
        state.commit(END_ACRONYMS)
        const { matched } = data.data || []
        state.commit(SET_ACRONYMS, { acronyms: matched })
        resolve(data.matched)
      })
    })
  }
}

const mutations = {
  [SET_ACRONYMS] (state, { acronyms }) {
    state.acronyms = acronyms
  },
  [BEGIN_ACRONYMS] (state) {
    state.loading = true
  },
  [END_ACRONYMS] (state) {
    state.loading = false
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
