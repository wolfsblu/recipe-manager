import {writable} from "svelte/store";

export interface Toast {
    id: number
    message: string
    timeout: number
    visible: boolean
    type: 'error' | 'info' | 'success' | 'warning',
    group: string
}

export type NewToast = Pick<Toast, 'message'> & {
    type?: Toast['type']
    group?: Toast['group']
}

const getDefaultToast = (): Toast => ({
    id: Math.floor(Math.random() * 1000),
    message: "",
    timeout: 10000,
    visible: true,
    type: 'info',
    group: 'default',
})

export const toasts = writable<Toast[]>([]);

export const addToast = (toast: NewToast) => {
    const newToast = Object.assign(getDefaultToast(), toast)
    toasts.update((all) => [newToast, ...all.filter(t => t.group !== toast.group)])
    setTimeout(() => dismissToast(newToast.id), newToast.timeout)
}

export const dismissToast = (id: number) => {
    toasts.update((all) => {
        const idx = all.findIndex(t => t.id === id)
        if (idx !== -1) {
            all[idx].visible = false;
        }
        return [...all];
    })
}