import ResponseError from '@/error/ResponseError'

export const buildError = function (prefix, error, suffix = '') {
  if (!process.env.VUE_APP_DEBUG) {
    prefix = ''
    suffix = ''
  } else {
    prefix = prefix + ' '
    suffix = ' ' + suffix
  }
  try {
    return new ResponseError(`${prefix}${error.response.data.message}${suffix}`, error.response)
  } catch {
    console.log('no status', error)
    return new ResponseError(`${prefix}${error.message}${suffix}`)
  }
}
