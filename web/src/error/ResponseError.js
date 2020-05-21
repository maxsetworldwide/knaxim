
export default class extends Error {
  constructor (msg, response = {}) {
    if (process.env.VUE_APP_DEBUG) {
      msg = `${response.status || 'No Status Code'} ${msg}`
    }
    super(msg)
    this.name = process.env.VUE_APP_DEBUG ? 'ResponseError' : 'Error'
    this.response = response
  }

  addDebug (msg) {
    if (process.env.VUE_APP_DEBUG) {
      this.message = `${this.message} ${msg}`
    }
    return this
  }
}
