import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const PermissionService = {
  permissions ({ id }) {
    return ApiService.get(`perm/file/${id}`).catch(error => {
      throw buildError('PermissionService.permissions', error)
    })
  },

  share ({ id, targets }) {
    return ApiService.post(`perm/file/${id}`, { 'id': targets }).catch(error => {
      throw buildError('PermissionService.share', error)
    })
  },

  stopShare ({ id, targets }) {
    return ApiService.delete(`perm/file/${id}`, { 'id': targets }).catch(error => {
      throw buildError('PermissionService.stopShare', error)
    })
  },
  makePublic ({ fid }) {
    return ApiService.post(`perm/file/${fid}/public`).catch(error => {
      throw buildError('PermissionService.makePublic', error)
    })
  },

  stopPublic ({ fid }) {
    return ApiService.delete(`perm/file/${fid}/public`).catch(error => {
      throw buildError('PermissionService.stopPublic', error)
    })
  }
}

export default PermissionService
