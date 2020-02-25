import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const GroupService = {
  create ({ name }) {
    return ApiService.put(`group`, { 'newname': name }).catch(error => {
      throw buildError('GroupService.create', error)
    })
  },

  // get associated groups, provide gid to get groups gid belongs to
  associated ({ gid }) {
    return ApiService.get(`group/options${gid ? '/' + gid : ''}`).catch(error => {
      throw buildError('GroupService.associated', error)
    })
  },

  info ({ gid }) {
    return ApiService.get(`group/${gid}`).catch(error => {
      throw buildError('GroupService.info', error)
    })
  },

  add ({ gid, target }) {
    return ApiService.post(`group/${gid}/member`, { 'id': target }).catch(error => {
      throw buildError('GroupService.add', error)
    })
  },

  remove ({ gid, target }) {
    return ApiService.delete(`group/${gid}/member`, { 'id': target }).catch(error => {
      throw buildError('GroupService.remove', error)
    })
  },

  lookup ({ name }) {
    return ApiService.get(`group/name/${name}`).catch(error => {
      throw buildError('GroupService.lookup', error)
    })
  }
}

export default GroupService
