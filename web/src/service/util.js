import RequestError from '@/error/RequestError'

export const buildError = function (prefix, error, suffix = '') {
  try {
    return new RequestError(`${prefix}${error.response.data.message}${suffix}`, error.response.status)
  } catch {
    return new RequestError(`${prefix}${error.message}${suffix}`)
  }
}
