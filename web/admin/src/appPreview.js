import { api } from './api';
export async function openAppPreview(id) {
    // Open a tab while still inside the click handler so popup blockers do not
    // prevent the preview after the API request completes.
    const previewWindow = window.open('about:blank', '_blank');
    try {
        const result = await api.post(`/apps/${id}/preview`);
        const target = new URL(result.url, window.location.origin);
        if (target.origin !== window.location.origin) {
            throw new Error('预览地址无效');
        }
        if (previewWindow && !previewWindow.closed) {
            previewWindow.opener = null;
            previewWindow.location.href = target.toString();
            return;
        }
        // If the browser blocked the new tab, keep the action useful by opening
        // the same preview in the current tab.
        window.location.href = target.toString();
    }
    catch (error) {
        previewWindow?.close();
        throw error;
    }
}
