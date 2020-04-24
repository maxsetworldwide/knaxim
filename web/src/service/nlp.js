import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const NLPService = {
  info ({ fid, category = 't', start = 0, end = 5 }) {
    return ApiService.get(`nlp/file/${fid}/${category}/${start}/${end}`).catch(err => {
      throw buildError('NLPService.info', err)
    })
  }
}

export default NLPService
