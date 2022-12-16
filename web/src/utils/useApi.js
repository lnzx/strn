import axios  from 'axios'

var api
export function createApi() {
    // Here we set the base URL for all requests made to the api
    api = axios.create({
        baseURL: import.meta.env.VITE_API_BASE_URL,
    })
    return api
}

export function useApi() {
    if (!api) {
        createApi()
    }
    return api
}



