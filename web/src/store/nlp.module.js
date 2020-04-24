import NLPService from '@/service/nlp'
import { NLP_DATA } from './actions.type'
import {
  LOADING_NLP,
  SET_NLP,
  PUSH_ERROR
} from './mutations.type'

function cat2state (cat) {
  switch (cat) {
    case 't':
    case 'topic':
      return 'topics'
    case 'a':
    case 'action':
      return 'actions'
    case 'r':
    case 'resource':
      return 'resources'
    case 'p':
    case 'process':
      return 'processes'
  }
}

const state = {
  topics: {}, // fid => []{word, count}
  actions: {},
  resources: {},
  processes: {},
  loading: 0
}

const getters = {
  nlpLoading (state) {
    return state.loading > 0
  },
  nlpTopics (state) {
    return state.topics
  },
  nlpActions (state) {
    return state.actions
  },
  nlpResources (state) {
    return state.resources
  },
  nlpProcesses (state) {
    return state.processes
  }
}

const mutations = {
  [LOADING_NLP] (state, delta) {
    state.loading += delta
  },
  [SET_NLP] (state, { fid, start, info, category }) {
    if (!state[cat2state(category)][fid]) {
      state[cat2state(category)][fid] = []
    }
    for (;state[cat2state(category)][fid].length < start + info.length;) {
      state[cat2state(category)][fid].push(null)
    }
    for (let i = 0; i < info.length; i++) {
      state[cat2state(category)][fid][start + i] = info[i]
    }
  }
}

const actions = {
  async [NLP_DATA] ({ commit, state }, { fid, category, start, end, overwrite }) {
    overwrite = overwrite ||
    !state[cat2state(category)][fid] ||
    state[cat2state(category)][fid].length < end
    for (let i = start; !overwrite && end > i; i++) {
      overwrite = !state[cat2state(category)][fid][i]
    }
    let result = null
    if (overwrite) {
      try {
        commit(LOADING_NLP, 1)
        let data = await NLPService.info({
          fid,
          category,
          start,
          end
        }).then(res => res.data)
        commit(SET_NLP, {
          fid,
          start,
          info: data.info,
          category
        })
        result = data.info
      } catch (e) {
        commit(PUSH_ERROR, new Error(`NLP_DATA: ${e}`))
      } finally {
        commit(LOADING_NLP, -1)
      }
    } else {
      result = state[cat2state(category)][fid].slice(start, end)
    }
    return result
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
