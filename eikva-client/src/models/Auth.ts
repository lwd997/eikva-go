export interface TokenResponse {
    access_token: string;
    refresh_token: string;
}

export interface WhoAmIResponse {
    login: string;
    uuid: string;
}

export interface ErrorResponse {
    error: string;
}

export interface ErrorValidResponse {
    form_errors: Array<ErrorValidField>
}

export interface ErrorValidField {
    field: string,
    error: string
}

export interface LogoutResponse {
    status: boolean;
}

