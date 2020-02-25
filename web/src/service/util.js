export const buildError = function (prefix, error, suffix = '') {
  try {
    return new Error(`${prefix}(${error.response.status})${error.response.data.message}${suffix}`)
  } catch {
    return new Error(`${prefix}${error.message}${suffix}`)
  }
}
