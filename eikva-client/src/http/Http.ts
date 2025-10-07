import { type TokenResponse } from "../models/Auth";
import { createErrorToast } from "../Toast";

interface HttpRequestInit extends Omit<RequestInit, 'headers' | 'body'> {
    headers?: Record<string, string>;
    body?: Record<string, unknown> | string | FormData;
}

type HttpBadStatus = 400 | 401 | 403 | 404 | 500;

interface HttpResponseOk<T extends object> {
    body: T;
    status: 200;
}

interface HttpResponseFail {
    body: null
    status: Exclude<HttpBadStatus, 200>;
}

type HttpResponse<T extends object> = HttpResponseOk<T> | HttpResponseFail;

interface HttpConstructor {
    baseUrl: string;
    accessToken: string | null;
    refreshToken: string | null;
    onUpdateTokenFail?: (h: Http) => void;
}


export class Http {
    public readonly baseUrl: string;
    private isRefreshing: boolean;
    private refreshPromise: Promise<void> | null;

    public accessToken: string | null;
    public refreshToken: string | null;
    public onUpdateTokenFail: ((h: Http) => void) | null;

    constructor(config: HttpConstructor) {
        this.baseUrl = config.baseUrl;
        this.isRefreshing = false;
        this.refreshPromise = null;
        this.accessToken = config.accessToken;
        this.refreshToken = config.refreshToken;
        this.onUpdateTokenFail = config.onUpdateTokenFail || null;
    }

    updateTokens({ refresh_token, access_token }: TokenResponse) {
        this.accessToken = access_token;
        this.refreshToken = refresh_token;
        localStorage.setItem("access_token", access_token);
        localStorage.setItem("refresh_token", refresh_token);
    }

    deleteTokens() {
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
    }

    private async waitForRefresh(): Promise<void> {
        if (this.isRefreshing && this.refreshPromise) {
            await this.refreshPromise;
        }
    }

    private async refreshTokens(): Promise<void> {
        if (this.isRefreshing) {
            return this.waitForRefresh();
        }

        this.isRefreshing = true;
        this.refreshPromise = (async () => {
            try {
                const refreshResponse = await this.request<TokenResponse>(
                    "/auth/update-tokens",
                    {
                        method: "POST",
                        body: { refresh_token: this.refreshToken },
                    }
                );

                if (refreshResponse.status === 200) {
                    this.updateTokens(refreshResponse.body);
                } else {
                    throw new Error("Refresh failed");
                }
            } catch (err) {
                this.deleteTokens();
                this.onUpdateTokenFail?.(this);
                throw err;
            } finally {
                this.isRefreshing = false;
                this.refreshPromise = null;
            }
        })();

        await this.refreshPromise;
    }

    async request<T extends object>(path: string, options: HttpRequestInit = {}): Promise<HttpResponse<T>> {
        if (!options.headers) {
            options.headers = {};
        }

        if (!options.headers["Content-Type"] && !(options.body instanceof FormData)) {
            options.headers["Content-Type"] = "application/json";
        }

        if (this.accessToken) {
            options.headers["Authorization"] = "Bearer " + this.accessToken;
        }

        if (options.body && typeof options.body !== "string" && !(options.body instanceof FormData)) {
            options.body = JSON.stringify(options.body);
        }

        await this.waitForRefresh();

        const resp = await fetch(this.baseUrl + path, options as RequestInit);

        if (resp.status === 401 && path !== "/auth/update-tokens") {
            try {
                await this.refreshTokens();
            } catch {
                return { status: 401, body: null };
            }

            options.headers["Authorization"] = "Bearer " + this.accessToken;
            const retryResp = await fetch(this.baseUrl + path, options as RequestInit);

            if (!retryResp.headers.get("Content-Type")?.includes("application/json")) {
                return { status: 500, body: null };
            }

            const retryBody = await retryResp.json();
            return { status: retryResp.status as HttpBadStatus | 200, body: retryBody };
        }

        if (!resp.headers.get("Content-Type")?.includes("application/json")) {
            return { status: 500, body: null };
        }

        const body = await resp.json();

        if (typeof body.error === "string") {
            createErrorToast(body.error);
        } else if (Array.isArray(body.form_errors)) {
            for (const fe of body.form_errors) {
                createErrorToast(fe.error);
            }
        }

        return {
            status: resp.status as HttpBadStatus | 200,
            body,
        };
    }
}

