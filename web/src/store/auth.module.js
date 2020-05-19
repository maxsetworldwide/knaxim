import UserService from '@/service/user'
import {
  LOGIN,
  AFTER_LOGIN,
  LOGOUT,
  REGISTER,
  GET_USER,
  CHANGE_PASSWORD,
  SEND_RESET_REQUEST,
  RESET_PASSWORD
} from './actions.type'
import {
  SET_USER,
  PURGE_AUTH,
  PROCESS_SERVER_STATE,
  AUTH_LOADING,
  PUSH_ERROR
} from './mutations.type'

const state = {
  user: null,
  loading: 0
}

const getters = {
  currentUser (state) {
    return state.user || {}
  },
  isAuthenticated (state) {
    return !!state.user
  },
  authLoading ({ loading }) {
    return loading > 0
  }
}

const actions = {
  /**
   * login - Sign in.
   *
   * @param {object} context  State
   * @param {object} credentials  Login crdentials
   * @param {string} credentials.login  Username
   * @param {string} credentials.password  Password
   * @return {Promise}
   */
  async [LOGIN] ({ commit, dispatch }, { login, password }) {
    commit(PURGE_AUTH)
    commit(AUTH_LOADING, 1)
    let out = null
    try {
      let res = await UserService.login({
        name: login,
        pass: password
      })
      dispatch(AFTER_LOGIN)
      out = res.data
    } catch (err) {
      commit(PUSH_ERROR, err.addDebug('action LOGIN'))
      throw err
    } finally {
      commit(AUTH_LOADING, -1)
    }
    return out
  },

  /**
   * Distroy login credentials on the server and client.
   *
   * @param {object} context  State
   * @return {Promise}
   */
  async [LOGOUT] (context) {
    context.commit(PURGE_AUTH)
    context.commit(AUTH_LOADING, 1)
    await UserService.logout()
      .catch((err) => {
        context.commit(PUSH_ERROR, err.addDebug('action LOGOUT'))
      })
      .finally(() => {
        context.commit(AUTH_LOADING, -1)
      })
  },

  /**
   * register - Create a new user.
   *
   * @param {object} context  State
   * @param {object} credentials  User account info
   * @param {string} credentails.email  An email address
   * @param {string} credentials.login  Login id
   * @param {string} credentials.password  Password
   * @return {Promise}
   */
  [REGISTER] ({ commit }, { email, login, password }) {
    commit(AUTH_LOADING, 1)
    return UserService.create({
      email,
      name: login,
      password
    }).then(({ data }) => data)
      .finally(() => commit(AUTH_LOADING, -1))
  },

  [CHANGE_PASSWORD] ({ commit, dispatch }, { oldpass, newpass }) {
    commit(AUTH_LOADING, 1)
    return UserService.changePassword({ oldpass, newpass })
      .then(() => dispatch(LOGOUT))
      .catch(err => commit(PUSH_ERROR, err.addDebug('action CHANGE_PASSWORD'))
      .finally(() => commit(AUTH_LOADING, -1))
  },

  [SEND_RESET_REQUEST] ({ commit }, { name }) {
    commit(AUTH_LOADING, 1)
    return UserService.requestReset({ name })
      .catch(err => commit(PUSH_ERROR, err.addDebug('action SEND_RESET_REQUEST'))
      .finally(() => commit(AUTH_LOADING, -1))
  },

  [RESET_PASSWORD] ({ commit }, { passkey, newpass }) {
    commit(AUTH_LOADING, 1)
    return UserService.resetPass({ passkey, newpass })
      .catch(err => commit(PUSH_ERROR, err.addDebug('action RESET_PASSWORD'))
      .finally(() => commit(AUTH_LOADING, -1))
  },

  /**
   * getUser - Get user data.
   *
   * @return {Promise}
   */
  [GET_USER] (context, { quiet } = { quiet: true }) {
    context.commit(AUTH_LOADING, 1)
    return UserService.info({})
      .then(({ data }) => {
        context.commit(SET_USER, data)
        return data
      })
      .catch((err) => {
        quiet || context.commit(PUSH_ERROR, `GET_USER: ${err}`)
        throw err
      })
      .finally(() => context.commit(AUTH_LOADING, -1))
  }
}

const mutations = {
  [SET_USER] (state, user) {
    state.user = user
  },

  [PURGE_AUTH] (state) {
    state.user = null
  },

  [PROCESS_SERVER_STATE] (state, { user }) {
    state.user = {
      id: user.id,
      name: user.name,
      data: user.data
    }
  },

  [AUTH_LOADING] (state, delta) {
    state.loading += delta
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
