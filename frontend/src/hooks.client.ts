import type {ClientInit} from "@sveltejs/kit";
import {Configuration, DefaultConfig} from "$lib/api";
import {PUBLIC_BACKEND_URL} from "$env/static/public";

export const init: ClientInit = async () => {
    DefaultConfig.config = new Configuration({
        basePath: PUBLIC_BACKEND_URL || '/',
        credentials: 'include',
    })
};