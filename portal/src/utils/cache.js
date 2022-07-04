export function setCache(k, v) {
  const content = JSON.stringify({
    type: typeof v,
    value: v
  })
  window.localStorage.setItem(k, content)
}

export function getCache(k, defaultValue) {
  const content = window.localStorage.getItem(k)
  try {
    const obj = JSON.parse(content)
    return obj.value
  } catch (e) {
    return defaultValue === undefined ? content : defaultValue
  }
}
