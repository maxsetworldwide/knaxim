import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import FileViewer from '@/components/file-viewer'
import PDFDoc from '@/components/pdf/pdf-doc'
import ImageViewer from '@/components/image-viewer'
import TextViewer from '@/components/text-viewer'
import { flushPromises, TestStore } from './utils'
import { TOUCH } from '@/store/mutations.type'
import { GET_FILE } from '@/store/actions.type'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const localVue = createLocalVue()
localVue.use(Vuex)

const testID = 'testID'
const testName = 'testName'
const testCount = 22
const testStore = new TestStore({
  state: {
    [testID]: {
      id: testID,
      name: testName,
      count: testCount
    }
  },
  actions: {
    [GET_FILE] (context, { id }) {
      return context.state[id]
    }
  },
  mutations: {
    [TOUCH] () {}
  }
})

const shallowMountFa = (
  options = { props: {}, methods: {}, computed: {}, store: null }
) => {
  let store = options.store || testStore.createStore()
  return shallowMount(FileViewer, {
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
      ...options.computed
    }
  })
}

describe('FileViewer', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(FileViewer)).toBe(true)
  })
  it('uses alternate views and passes props to them', async () => {
    const expectViews = function (views, expected) {
      for (const view in views) {
        if (expected[view]) {
          expect(views[view].exists()).toBeTrue()
        } else {
          expect(views[view].exists()).toBeFalse()
        }
      }
    }
    const findViews = function (wrapper) {
      const pdf = wrapper.find(PDFDoc)
      const image = wrapper.find(ImageViewer)
      const text = wrapper.find(TextViewer)
      return {
        pdf,
        image,
        text
      }
    }

    const wrapper = shallowMountFa()

    // order matters, semi dependent on implementation
    let views = findViews(wrapper)
    expectViews(views, { pdf: true, image: false, text: false })
    const pdfdoc = wrapper.find(PDFDoc)
    expect(pdfdoc.props()).toEqual({ fileID: testID, acr: undefined })
    pdfdoc.vm.$emit('no-view')
    await flushPromises()

    views = findViews(wrapper)
    expectViews(views, { pdf: false, image: true, text: false })
    const image = wrapper.find(ImageViewer)
    expect(image.props()).toEqual({ id: testID })
    image.vm.$emit('no-image')
    await flushPromises()

    views = findViews(wrapper)
    expectViews(views, { pdf: false, image: false, text: true })
    const text = wrapper.find(TextViewer)
    expect(text.props()).toEqual({
      fileName: testName,
      finalPage: testCount,
      acr: undefined
    })
  })
  it('touches file', () => {
    const store = testStore.createStore()
    spyOn(store, 'commit')
    shallowMountFa({ store })
    const touchCommits = store.commit.calls.allArgs().filter((args) => {
      return args[0] === TOUCH
    })
    expect(touchCommits.length).toBeGreaterThan(0)
    touchCommits.forEach((call) => {
      expect(call[1]).toEqual(testID)
    })
  })
})
