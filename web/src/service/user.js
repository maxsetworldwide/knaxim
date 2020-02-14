import ApiService from '@/service/api'

const UserService = {
  create ({ name, password, email }) {
    return ApiService.put(`user`, { 'name': name, 'pass': password, 'email': email }).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  },

  createAdmin ({ name, password, email, adminKey }) {
    return ApiService.put(`user/admin`, { 'name': name, 'pass': password, 'email': email, 'adminKey': adminKey }).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  },

  info ({ id }) {
    if (id === undefined) {
      return ApiService.get(`user`).catch(error => {
        throw new Error(`UserService ${error}`)
      })
    } else {
      return ApiService.query(`user`, { 'id': id }).catch(error => {
        throw new Error(`UserService ${error}`)
      })
    }
  },

  lookup ({ name }) {
    return ApiService.get(`user/name/${name}`).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  },

  login ({ name, pass }) {
    return ApiService.post(`user/login`, { 'name': name, 'pass': pass }).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  },

  logout () {
    return ApiService.delete(`user`).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  },

  completeProfile () {
    return ApiService.get(`user/complete`).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  },

  changePassword ({ oldpass, newpass }) {
    return ApiService.post(`user/pass`, { 'oldpass': oldpass, 'newpass': newpass }).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  },

  data () {
    return ApiService.get(`user/data`).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  },

  requestReset ({ name }) {
    return ApiService.put('user/reset', { name }).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  },

  resetPass ({ key, newpass }) {
    return ApiService.post('user/reset', { key, newpass }).catch(error => {
      throw new Error(`UserService ${error}`)
    })
  }
}

export default UserService
