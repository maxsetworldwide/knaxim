
export default class extends Error {
  constructor (msg, status = 0) {
    if (process.env.VUE_APP_DEBUG) {
      msg = `${status || 'No Status Code'} ${msg}`
    }
    super(msg)
    this.name = process.env.VUE_APP_DEBUG ? 'RequestError' : 'Error'
    this.status = status
  }

  addDebug (msg) {
    if (process.env.VUE_APP_DEBUG) {
      this.message = `${this.message} ${msg}`
    }
    return this
  }
}
