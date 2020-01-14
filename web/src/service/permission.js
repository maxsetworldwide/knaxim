import ApiService from '@/service/api'

const PermissionService = {
  permissions ({ id }) {
    return ApiService.get(`perm/file/${id}`).catch(error => {
      throw new Error(`FileService ${error}`)
    })
  },

  share ({ id, targets }) {
    return ApiService.post(`perm/file/${id}`, { 'id': targets }).catch(error => {
      throw new Error(`FileService ${error}`)
    })
  },

  stopShare ({ id, targets }) {
    return ApiService.delete(`perm/file/${id}`, { 'id': targets }).catch(error => {
      throw new Error(`FileService ${error}`)
    })
  },
  makePublic ({ fid }) {
    return ApiService.post(`perm/file/${fid}/public`).catch(error => {
      throw new Error(`AdminService ${error}`)
    })
  },

  stopPublic ({ fid }) {
    return ApiService.delete(`perm/file/${fid}/public`).catch(error => {
      throw new Error(`AdminService ${error}`)
    })
  }
}

export default PermissionService
