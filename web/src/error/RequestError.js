
export default class extends Error {
  constructor(msg, status=0) {
    super(msg)
    this.status = status
    this.debug = []
  }

  addDebug(msg) {
    this.debug.push(msg)
    return this
  }

  get message() {
    if (process.env.VUE_APP_DEBUG) {
      return [status||'', super.message, this.debug...].join('\n')
    } else {
      return super.message
    }
  }
}
