import {goto} from "$app/navigation";
import {OpenIdAuthServiceApi} from "$lib/api";
import toasts from "$lib/state/Toast.svelte";
import {m} from '$lib/paraglide/messages.js';

const EmptyBody = {
    body: {
    }
};

class Auth {
    public loading: boolean = $state(true);
    private readonly authServiceApi: OpenIdAuthServiceApi;
    private userInfo: {
        subject?: string;
        email?: string;
    } = $state({
        subject: undefined,
        email: undefined
    });

    constructor() {
        this.authServiceApi = new OpenIdAuthServiceApi()

        this.initialize()
            .catch(async (error) => {
                console.error('Error during Auth initialization:', error);
                toasts.addToast(
                    m.generic_authentication_error(),
                    'error'
                )
                await this.toLogin();
            });
    }

    public getUserInfo(): typeof this.userInfo {
        return this.userInfo;
    }

    async initialize() {
        const urlParams = new URLSearchParams(window.location.search);
        const code = urlParams.get('code');
        const state = urlParams.get('state');
        if (code && state) {
            return this.exchangeCode(code, state);
        }

        try {
            this.userInfo = await this.authServiceApi.openIdAuthServiceGetUserInfo(EmptyBody);
            this.loading = false;

            if (window.location.pathname.startsWith('/login')) {
                return this.toDashboard();
            }

            return;
        }
        catch (error) {
            console.error('Error fetching user info:', error);
        }

        await this.beginAuth();
    }

    async beginAuth(): Promise<void> {
        const authUrl = await this.authServiceApi.openIdAuthServiceBeginAuth(EmptyBody)
        console.log('Redirecting to auth URL:', authUrl.url);
        if (!authUrl.url) {
            console.error('Auth URL is empty');
            toasts.addToast(
                m.generic_authentication_error(),
                'error'
            );
            return;
        }

        window.location.href = authUrl.url.replace('http://dex', 'http://localhost');
    }

    async exchangeCode(code: string, state: string): Promise<void> {
        const url = new URL(window.location.href);
        url.searchParams.delete('code');
        url.searchParams.delete('state');
        await goto(url, {
            replaceState: true,
        })

        try {
            await this.authServiceApi.openIdAuthServiceExchangeCode({
                body: {
                    state,
                    code
                }
            })

            return this.initialize();
        }
        catch (error) {
            console.error('Error exchanging code:', error);
            toasts.addToast(
                m.generic_authentication_error(),
                'error'
            );
            return;
        }
    }

    async logout(): Promise<void> {
        try {
            await this.authServiceApi.openIdAuthServiceLogout(EmptyBody);
        } catch (error) {
            console.error('Error during logout:', error);
            toasts.addToast(
                m.generic_authentication_error(),
                'error'
            );
        }

        this.userInfo = {
            subject: undefined,
            email: undefined
        }
        await this.toLogin();
    }

    async toLogin(): Promise<void> {
        await goto(
            '/login',
            {
                replaceState: true,
            }
        );
    }

    async toDashboard(): Promise<void> {
        await goto(
            '/',
            {
                replaceState: true,
            }
        );
    }
}

const auth = new Auth();

export default auth;