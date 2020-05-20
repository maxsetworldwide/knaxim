import RequestError from '@/error/RequestError'

export const buildError = function (prefix, error, suffix = '') {
  try {
    console.log('building error')
    console.log(error)
    console.log(error.response)
    return new RequestError(`${prefix} ${error.response.data.message} ${suffix}`, error.response.status)
  } catch {
    console.log('no status', error)
    return new RequestError(`${prefix} ${error.message} ${suffix}`)
  }
}
