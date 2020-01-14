import ApiService from '@/service/api'

const AcronymService = {
  get ({ acronym }) {
    return ApiService.get(`acronym/${acronym}`).catch(error => {
      throw new Error(`AcronymService ${error}`)
    })
  }
}

export default AcronymService
