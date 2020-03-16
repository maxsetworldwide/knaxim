
// helper for testing action with expected mutations
export const testAction = (action, payload, expectedMutations, done, context = {}) => {
  let count = 0

  // mock commit
  const commit = (type, payload) => {
    const mutation = expectedMutations[count]

    try {
      expect(type).to.equal(mutation.type)
      if (payload) {
        expect(payload).to.deep.equal(mutation.payload)
      }
    } catch (error) {
      done(error)
    }

    count++
    if (count >= expectedMutations.length) {
      done()
    }
  }

  // call the action with mocked store and arguments
  let result = action({ ...context, commit }, payload)
  if (result.finally) { // if action returned a promise
    result.finally(() => {
      if (expectedMutations.length === 0) {
        expect(count).to.equal(0)
        done()
      }
    })
  } else {
    if (expectedMutations.length === 0) {
      expect(count).to.equal(0)
      done()
    }
  }
}
