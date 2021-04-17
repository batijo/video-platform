export const toCamelCase = (str: string): string => {
  let strArr = str.split('_')
  if (strArr.length === 1) return str.toLowerCase()

  return strArr.slice(1).reduce((acc, n) => {
    return acc += n[0].toUpperCase() + n.slice(1)
  }, strArr[0].toLowerCase())
}

export const toCamelCaseObj = (obj: any): any => {
  for (let key in obj) {
    let newKey = toCamelCase(key)

    Object.assign(obj, { [newKey]: obj[key] })[key]

    if (typeof obj[newKey] === 'object') {
      obj[newKey] = toCamelCaseObj(obj[newKey])
    }

    if (key !== newKey) delete obj[key]
  }

  return obj
}