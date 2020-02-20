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
  SET_ERROR,
  PROCESS_SERVER_STATE,
  AUTH_LOADING,
  PUSH_ERROR
} from './mutations.type'

const state = {
  errors: null,
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
      if (res.data.message === 'Not Found') {
        throw new Error('Failed to login')
      }
      dispatch(AFTER_LOGIN)
      out = res.data
    } catch (err) {
      // TODO: handle Error
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
  [LOGOUT] (context) {
    context.commit(PURGE_AUTH)
    context.commit(AUTH_LOADING, 1)
    UserService.logout().then(({ data }) => {
    }).catch(({ response }) => {
      context.commit(SET_ERROR, response.data)
    }).finally(() => {
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
    return UserService.changePassword({ oldpass, newpass }).then(() => dispatch(LOGOUT)).finally(() => commit(AUTH_LOADING, -1))
  },

  [SEND_RESET_REQUEST] ({ commit }, { name }) {
    commit(AUTH_LOADING, 1)
    return UserService.requestReset({ name }).finally(() => commit(AUTH_LOADING, -1))
  },

  [RESET_PASSWORD] ({ commit }, { passkey, newpass }) {
    commit(AUTH_LOADING, 1)
    return UserService.resetPass({ passkey, newpass }).finally(() => commit(AUTH_LOADING, -1))
  },

  /**
   * getUser - Get user data.
   *
   * @return {Promise}
   */
  [GET_USER] (context) {
    context.commit(AUTH_LOADING, 1)
    return (new Promise((resolve, reject) => {
      UserService.info({})
        .then(({ data }) => {
          if (data.message === 'login') {
            context.commit(SET_ERROR, data.message)
            reject(data)
          } else {
            context.commit(SET_USER, data)
            resolve(data)
          }
        }).catch(({ response }) => {
          if (response && response.message) {
            context.commit(SET_ERROR, response.message)
            reject(new Error(response.message))
          } else {
            context.commit(SET_ERROR, 'unable to get user')
            reject(new Error('unable to get user'))
          }
        })
    })).finally(() => context.commit(AUTH_LOADING, -1))
  }
}

const mutations = {
  [SET_ERROR] (state, error) {
    state.errors = error
  },

  [SET_USER] (state, user) {
    state.user = user
    state.errors = null
  },

  [PURGE_AUTH] (state) {
    state.user = null
    state.errors = null
  },

  [PROCESS_SERVER_STATE] (state, { user }) {
    state.user = {
      id: user.id,
      name: user.name,
      data: user.data
    }
    state.errors = null
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
