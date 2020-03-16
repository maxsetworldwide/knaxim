
// helper for testing action with expected mutations
export const testAction = (action, payload, expectedMutations, done, context = {}) => {
  let count = 0

  // mock commit
  const commit = (type, payload) => {
    const mutation = expectedMutations[count]

    try {
      expect(type).toBe(mutation.type)
      if (payload) {
        expect(payload).toEqual(mutation.payload)
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
  if (result && result.finally) { // if action returned a promise
    result.finally(() => {
      if (expectedMutations.length === 0) {
        expect(count).toBe(0)
        done()
      }
    })
  } else {
    if (expectedMutations.length === 0) {
      expect(count).toBe(0)
      done()
    }
  }
}
