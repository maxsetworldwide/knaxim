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
  SET_ERROR
} from './mutations.type'

const state = {
  errors: null,
  user: {},
  isAuthenticated: false
}

const getters = {
  currentUser (state) {
    return state.user
  },
  isAuthenticated (state) {
    return state.isAuthenticated
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
    return new Promise((resolve, reject) => {
      UserService.login({
        name: credentials.login,
        pass: credentials.password
      }).then(({ data }) => {
        if (data.message === 'Not Found') {
          context.commit(SET_ERROR, data)
          reject(data)
        } else {
          context.dispatch(GET_USER, {})
          context.dispatch(AFTER_LOGIN, {})
          resolve(data)
        }
      }).catch(({ response }) => {
        context.commit(SET_ERROR, response.data)
        reject(response)
      })
    })
  },

  /**
   * Distroy login credentials on the server and client.
   *
   * @param {object} context  State
   * @return {Promise}
   */
  [LOGOUT] (context) {
    context.commit(PURGE_AUTH)
    UserService.logout().then(({ data }) => {
    }).catch(({ response }) => {
      context.commit(SET_ERROR, response.data)
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
    return new Promise((resolve, reject) => {
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
    })
  },

  /**
   * getUser - Get user data.
   *
   * @return {Promise}
   */
  [GET_USER] (context) {
    return new Promise((resolve, reject) => {
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
          context.commit(SET_ERROR, response.data.errors)
        })
    })
  }
}

const mutations = {
  [SET_ERROR] (state, error) {
    state.errors = error
  },

  [SET_USER] (state, user) {
    state.isAuthenticated = true
    state.user = user
    state.errors = {}
  },

  [PURGE_AUTH] (state) {
    state.isAuthenticated = false
    state.user = {}
    state.errors = {}
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
