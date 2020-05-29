import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import ImageViewer from '@/components/image-viewer'
import FileActions from '@/components/file-actions'
import { TestStore } from './utils'

const localVue = createLocalVue()
localVue.use(Vuex)

const testID = 'testID'
const testName = 'testName'
const testSrc = 'testSrc'
const testFile = { name: testName }
const testStore = new TestStore({
  state: {
    fileset: {
      [testID]: testFile
    }
  },
  getters: {
    populateFiles (state) {
      return function (id) {
        return state.fileset[id]
      }
    }
  }
})

const shallowMountFa = (
  options = { props: {}, methods: {}, computed: {}, store: null }
) => {
  let store = options.store || testStore.createStore()
  return shallowMount(ImageViewer, {
    localVue,
    store,
    propsData: {
      id: testID,
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      srcURL () {
        return testSrc
      },
      ...options.computed
    }
  })
}

describe('ImageViewer', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(ImageViewer)).toBe(true)
  })
  it('renders title', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.text()).toContain(testName)
  })
  it('renders image with correct src', () => {
    const wrapper = shallowMountFa()
    const img = wrapper.find('img')
    expect(img.attributes('src')).toEqual(testSrc)
  })
  it('propogates img error as "no-image"', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.emitted('no-image')).toBeUndefined()
    wrapper.find('img').trigger('error')
    expect(wrapper.emitted('no-image').length).toEqual(1)
  })
  it('uses file-actions correctly', () => {
    const wrapper = shallowMountFa()
    const fileActions = wrapper.find(FileActions)
    expect(fileActions.props()).toEqual({
      checkedFiles: [testFile],
      singleFile: true,
      disableDownloadPDF: true
    })
  })
})
