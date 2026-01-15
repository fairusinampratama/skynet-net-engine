import axios from "axios";

export const api = axios.create({
    baseURL: "/api/v1",
    headers: {
        "Content-Type": "application/json",
    },
});

// Add interceptor to inject API Key
api.interceptors.request.use((config) => {
    config.headers["X-App-Key"] = "netengine_secret_key_123";
    return config;
});
