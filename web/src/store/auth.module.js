import UserService from '@/service/user'
import {
  LOGIN,
  AFTER_LOGIN,
  LOGOUT,
  REGISTER,
  GET_USER
} from './actions.type'
import {
  SET_USER,
  PURGE_AUTH,
  SET_ERROR,
  PROCESS_SERVER_STATE,
  AUTH_LOADING
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
  [LOGIN] (context, credentials) {
    context.commit(PURGE_AUTH)
    context.commit(AUTH_LOADING, 1)
    return (new Promise((resolve, reject) => {
      UserService.login({
        name: credentials.login,
        pass: credentials.password
      }).then(({ data }) => {
        if (data.message === 'Not Found') {
          context.commit(SET_ERROR, data)
          reject(data)
        } else {
          context.dispatch(AFTER_LOGIN, {})
          resolve(data)
        }
      }).catch(({ response }) => {
        context.commit(SET_ERROR, response.data)
        reject(response)
      })
    })).finally(() => context.commit(AUTH_LOADING, -1))
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
  [REGISTER] (context, credentials) {
    context.commit(AUTH_LOADING, 1)
    return (new Promise((resolve, reject) => {
      UserService.create({
        email: credentials.email,
        name: credentials.login,
        password: credentials.password
      }).then(({ data }) => {
        if (data.message === 'Name Already Taken' || data.message === 'Bad Request') {
          context.commit(SET_ERROR, data.message)
          reject(data)
        } else {
          resolve(data)
        }
      }).catch(({ response }) => {
        context.commit(SET_ERROR, response.data.errors)
        reject(response)
      })
    })).finally(() => context.commit(AUTH_LOADING, -1))
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
