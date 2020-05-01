import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const SearchService = {

  folderFiles ({ name, find, group }) {
    return ApiService.query(`dir/${name}/search`, { find, group }).catch(error => {
      throw buildError('SearchService.folderFiles', error)
    })
  },

  groupFiles ({ gid, find }) {
    return ApiService.query(`group/${gid}/search`, { find }).catch(error => {
      throw buildError('SearchService.groupFiles', error)
    })
  },

  publicFiles ({ find }) {
    return ApiService.query(`public/search`, { find }).catch(error => {
      throw buildError('SearchService.publicFiles', error)
    })
  },
  userFiles ({ find }) {
    return ApiService.query(`user/search`, { find }).catch(error => {
      throw buildError('SearchService.userFiles', error)
    })
  },
  FileTags ({ context, match }) {
    return ApiService.post(`search/tags`, {
      context,
      match,
      _sendJSON: true
    }).catch(error => {
      throw buildError('SearchService.FileTags', error)
    })
  },
  newOwnerContext (oid, limit) {
    let out = {
      type: 'owner',
      id: oid
    }
    if (limit) {
      out['only'] = limit
    }
    return out
  },
  newFileContext (fid) {
    let out = {
      type: 'file',
      id: fid
    }
    return out
  },
  newMatchCondition (find, type, regex = true, owner) {
    if (type) {
      return {
        tagtype: type,
        word: find,
        regex: regex,
        owner: owner
      }
    } else {
      return find
    }
  }
}

export default SearchService
