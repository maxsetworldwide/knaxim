import ApiService from '@/service/api'

const GroupService = {
  create ({ name }) {
    return ApiService.put(`group`, { 'newname': name }).catch(error => {
      throw new Error(`GroupService ${error}`)
    })
  },

  // get associated groups, provide gid to get groups gid belongs to
  associated ({ gid }) {
    return ApiService.get(`group/options${gid ? '/' + gid : ''}`).catch(error => {
      throw new Error(`GroupService ${error}`)
    })
  },

  info ({ gid }) {
    return ApiService.get(`group/${gid}`).catch(error => {
      throw new Error(`GroupService ${error}`)
    })
  },

  add ({ gid, target }) {
    return ApiService.post(`group/${gid}/member`, { 'id': target }).catch(error => {
      throw new Error(`GroupService ${error}`)
    })
  },

  remove ({ gid, target }) {
    return ApiService.delete(`group/${gid}/member`, { 'id': target }).catch(error => {
      throw new Error(`GroupService ${error}`)
    })
  },

  lookup ({ name }) {
    return ApiService.get(`group/name/${name}`).catch(error => {
      throw new Error(`GroupService ${error}`)
    })
  }
}

export default GroupService
