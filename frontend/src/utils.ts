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

export const toSnakeCase = (str: string) => str.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`)

export const toSnakeCaseObj = (obj: any) => {
  for (let key in obj) {
    let newKey = toSnakeCase(key)

    Object.assign(obj, { [newKey]: obj[key] })[key]

    if (typeof obj[newKey] === 'object') {
      obj[newKey] = toSnakeCaseObj(obj[newKey])
    }

    if (key !== newKey) delete obj[key]
  }

  return obj
}

// SSE To Do ...
export const sseProgress = (): any => {
  var source = new EventSource(`${window.origin}/api/sse/dashboard/`);
  console.log("Connection to /sse/dashboard established")
  var log = '';
  var currentmsg = '';

  source.onmessage = function (event) {
    if (!event.data.startsWith('<')) {
      localStorage.setItem('filename', event.data)
    } else if (event.data.startsWith('videouri')) {
      var temp = event.data;
      // manifestUri = temp.replace('videouri: ', '');
    } else if (event.data.indexOf('Error') > -1) {
      log += '<span class="error">' + event.data + '</span><br>';
    } else if (/^[\s\S]*<br>.*?Progress:.*?<br>$/.test(log) && event.data.includes('Progress:')) {
      log = log.replace(/^([\s\S]*<br>)(.*?Progress:.*?)(<br>)$/, `$1${event.data}$3`);
    } else if (event.data.indexOf('Transcoding complete') > -1 || event.data.indexOf('Transcoding parameters saved') > -1) {
      currentmsg = event.data;
      log += currentmsg + '<br>';
    } else {
      currentmsg = event.data;
      log += currentmsg + '<br>';
    }

    //document.getElementById('console').innerHTML = logg;
  };

  source.onerror = function (event) {
    console.log("Event Source failed.")
  }
}