// reactive Toast service class

class Toasts {
    toasts: Array<{ id: string; message: string; type: 'success' | 'error' | 'info' | 'warning' }> = $state([]);

    addToast(message: string, type: 'success' | 'error' | 'info' = 'info') {
        const id = crypto.randomUUID();
        this.toasts.push({ id, message, type });
        setTimeout(() => this.removeToast(id), 5000); // Auto-remove after 5 seconds
    }

    removeToast(id: string) {
        this.toasts = this.toasts.filter(toast => toast.id !== id);
    }
}

const toasts = new Toasts();

export default toasts;