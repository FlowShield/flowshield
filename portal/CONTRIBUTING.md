# CONTRIBUTING

First of all, thank you everyone who made pull requests for ZAPortal.

## 📦 Key Technical Skills

- [vue2](https://v2.vuejs.org/)
- [vuetify](https://vuetifyjs.com/)
- [vuex](https://vuex.vuejs.org/)
- [vue router](https://router.vuejs.org/)

## 📁 Directories

Some common usage of directories(files):

```shell
├── src
│   ├── App.vue
│   ├── api           # backend api request
│   ├── assets        # images assets
│   ├── components    # project scope components
│   ├── main.js       # the main process of APP, Vue instance global prototype, gloabl scss
│   ├── permission.js # login and auth logic, also include global top progress loading bar
│   ├── plugins       # vuetify plugin configurations
│   ├── router        # router defination
│   ├── store         # vuex
│   ├── styles 
│   │   ├── index.scss     # global css
│   │   └── variables.scss # override vuetify scss variables
│   ├── utils         # global helper functions
│   │   ├── event-bus.js      # eventbus
│   │   ├── request-helper.js # request interceptor
│   │   └── request.js        # wrapper of axios
│   └── views         # pages
├── tests
└── vue.config.js
```

## 💻 Quick Start

```bash
npm install
npm run serve
```

It servers on `http://localhost:8080` by default.

## 🎁 Read More

### Global message box

When error occurs, we want to alert a global message. You can do it simply by `EventBus`

```js
import { EventBus } from '@/utils/event-bus'

EventBus.$emit('app.message', 'Need login', 'warning')
```

### Global Backend API request handler

If the request does not reponse with `200` status code, the global error message will be trigger.

We handle the global response with axios interceptor, you can find more in file `src/utils/request-helper.js`.

```js
service.interceptors.response.use(
  error => {
    const statusCode = error.response.status
    
    if (statusCode === 401 && ['/login', '/'].includes(window.location.pathname)) { // ignore error message
      return Promise.reject(error)
    }

    EventBus.$emit('app.message', `[${statusCode}] ${error.message}`, 'error')
    return Promise.reject(error)
    
  }
)
```


### EventBus

| name            | description       |
|-                |-                  |
| `app.message`   | global message box|
| `app.loading`   | top page loading bar |


## Mock

Now you can use mock API to develop, which has been set as the default development API in `.env.development`.

```dotenv
VUE_APP_BASE_URL="https://531f6a00-a189-4209-a65c-f95e10e121cb.mock.pstmn.io/api/v1"
```
And if you have your own backend api, you can use environment variable to override mock api, eg:

```bash
export VUE_APP_BASE_URL=http://YOURHOST/api/v1 && npm run serve 
```
