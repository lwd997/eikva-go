import { appStore } from "../Storage";
import { Http } from "./Http";

export const http = new Http({
    accessToken: localStorage.getItem("access_token"),
    refreshToken: localStorage.getItem("refresh_token"),
    baseUrl: window.location.origin,
    onUpdateTokenFail: (h) => {
        h.deleteTokens();
        appStore.updateField('isAuthorized', false);
    }
});

