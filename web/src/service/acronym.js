import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const AcronymService = {
  get ({ acronym }) {
    return ApiService.get(`acronym/${acronym}`).catch(error => {
      throw buildError('AcronymService.get', error)
    })
  }
}

export default AcronymService
