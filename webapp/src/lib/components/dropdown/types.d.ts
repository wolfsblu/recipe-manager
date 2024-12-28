import type {Component} from "svelte";

export interface MenuItem {
    class?: string
    label: string
    href?: string
    icon?: Component
    onClick?: (e: MouseEvent) => void
}