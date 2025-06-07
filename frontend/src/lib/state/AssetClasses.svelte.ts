import {AssetClassServiceApi, type DigInvAssetClass} from "$lib/api";
import * as runtime from "$lib/api/runtime";
import toasts from "$lib/state/Toast.svelte";
import auth from "$lib/state/Auth.svelte";

abstract class ApiEndpoint {
    public loading: boolean = $state(false);
    protected abstract readonly apiInstance: runtime.BaseAPI;

    protected withLoading<T>(promise: Promise<T>): Promise<T> {
        this.loading = true;
        return promise.finally(() => {
            this.loading = false;
        });
    }

    protected async withErrorHandling<T>(
        promise: Promise<T>,
        errorMessage: string = 'An error occurred while processing your request.'
    ): Promise<T> {
        try {
            return await promise;
        } catch (error) {
            console.error('API Error:', error);
            toasts.addToast(
                errorMessage,
                'error'
            );

            throw error;
        }
    }

    // logout when 401 Unauthorized
    protected async handleUnauthorized<T>(promise: Promise<T>): Promise<T> {
        return promise.catch((error) => {
            if (error.status === 401) {
                console.warn('Unauthorized access, logging out...', error);
                auth.logout();
                return Promise.reject(new Error('Unauthorized'));
            }

            throw error;
        });
    }

    protected async withDefaults<T>(promise: Promise<T>): Promise<T> {
        return await this.withLoading(
            this.withErrorHandling(
                this.handleUnauthorized(promise)
            )
        );
    }
}

class AssetClassesBase extends ApiEndpoint {
    protected readonly apiInstance = new AssetClassServiceApi();
}

class AssetClasses extends AssetClassesBase {
    public assetClasses: Array<DigInvAssetClass> = $state([]);

    public async load(): Promise<void> {
        const response = await this.withDefaults(
            this.apiInstance.assetClassServiceGetAssetClasses({ body: {} })
        );
        this.assetClasses = response.classes || [];
    }
}

export const assetClasses: AssetClasses = new AssetClasses();

class AssetClassesCreate extends AssetClassesBase {
    public async create(assetClass: DigInvAssetClass): Promise<void> {
        const response = await this.withDefaults(
            this.apiInstance.assetClassServiceCreateAssetClass({ body: assetClass })
        );
        assetClasses.assetClasses.push(response);
    }
}

export const createAssetClasses: AssetClassesCreate = new AssetClassesCreate();
