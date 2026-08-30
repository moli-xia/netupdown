import axios from 'axios';
let accessToken = '';
let refreshing = null;
export const api = axios.create({ baseURL: '/api/admin', withCredentials: true, timeout: 30000 });
api.interceptors.request.use(c => { if (accessToken)
    c.headers.Authorization = `Bearer ${accessToken}`; return c; });
api.interceptors.response.use(r => r.data.data, async (err) => { const cfg = err.config; if (err.response?.data?.code === 10003 && !cfg.__retried) {
    cfg.__retried = true;
    refreshing ||= api.post('/auth/refresh').then((x) => { accessToken = x.access_token; return accessToken; }).finally(() => refreshing = null);
    const token = await refreshing;
    cfg.headers.Authorization = `Bearer ${token}`;
    return api.request(cfg);
} ; throw new Error(err.response?.data?.message || err.message); });
export function setToken(v) { accessToken = v; }
export async function bootstrapAuth() { try {
    const data = await api.post('/auth/refresh');
    setToken(data.access_token);
    return data.user;
}
catch {
    return null;
} }
