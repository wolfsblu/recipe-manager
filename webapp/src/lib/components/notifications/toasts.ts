import {writable} from "svelte/store";

export interface Toast {
    id: number
    message: string
    timeout: number
    visible: boolean
    progress: number
    showProgress: boolean
    type: 'error' | 'info' | 'success' | 'warning'
}

export type NewToast = Pick<Toast, 'message'> & {
    type?: Toast['type']
}

const getDefaultToast = (): Toast => ({
    id: Math.floor(Math.random() * 1000),
    message: "",
    progress: 100,
    timeout: 5000,
    visible: true,
    showProgress: true,
    type: 'info',
})

export const toasts = writable<Toast[]>([]);

export const dismissToast = (id: number) => {
    toasts.update((all) => {
        const idx = all.findIndex(t => t.id === id)
        if (idx !== -1) {
            all[idx].visible = false;
        }
        return [...all];
    })
}
export const addToast = (toast: NewToast) => {
    const newToast = Object.assign(getDefaultToast(), toast)
    toasts.update((all) => [newToast, ...all])
    setTimeout(() => dismissToast(newToast.id), newToast.timeout)
}