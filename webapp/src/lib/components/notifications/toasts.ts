import {writable} from "svelte/store";

export interface Toast {
    id: number
    createdAt: Date
    elapsed: number
    message: string
    timeout: number
    timeoutHandle: number
    visible: boolean
    pausedAt: Date
    progress: number
    type: 'error' | 'info' | 'success' | 'warning',
    group: string
}

export type NewToast = Pick<Toast, 'message'> & {
    type?: Toast['type']
    group?: Toast['group']
}

const getDefaultToast = (): Toast => ({
    id: Math.floor(Math.random() * 1000),
    createdAt: new Date(),
    message: "",
    progress: 100,
    timeout: 10000,
    elapsed: 0,
    timeoutHandle: 0,
    pausedAt: new Date(),
    visible: true,
    type: 'info',
    group: 'default',
})

export const toasts = writable<Toast[]>([]);

export const addToast = (toast: NewToast) => {
    const newToast = Object.assign(getDefaultToast(), toast)
    newToast.timeoutHandle = setTimeout(() => dismissToast(newToast.id), newToast.timeout)
    toasts.update((all) => [newToast, ...all.filter(t => t.type !== toast.type && t.group !== toast.group)])

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