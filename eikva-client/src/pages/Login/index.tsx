import { useState } from "react";
import { http } from "../../http";
import type { TokenResponse } from "../../models/Auth";
import { appStore } from "../../Storage";
import Button from "../../components/universal/Button/Button";
import "./Login.css";

export const Login = () => {
    const [formType, setFormType] = useState<"login" | "register">("login");
    const [login, setLogin] = useState("");
    const [password, setPassword] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const requestPath = formType === "login"
            ? "/auth/login"
            : "/auth/register"

        const response = await http.request<TokenResponse>(requestPath, {
            method: "POST",
            body: { login, password }
        });

        if (response.status === 200) {
            http.updateTokens(response.body);
            appStore.updateField('isAuthorized', true);
        }
    }

    const toggleFormType = () => {
        setFormType(formType === "login" ? "register" : "login");
    }

    let actionBtnTitle: string;
    let secondaryBtnTitle: string;

    switch(formType) {
        case "register":
            actionBtnTitle = "Зарегистроваться";
            secondaryBtnTitle = "Вход";
            break;
        case "login":
        default:
            actionBtnTitle = "Войти";
            secondaryBtnTitle = "Регистрация";
            break;
    }


    return (
        <form onSubmit={handleSubmit} className="display-flex width-100 height-100 align-items-center justify-content-center">
            <div className="login-form card display-flex flex-direction-column">
                <div className="logo">
                    <img src="/media/TestCraft.svg" />
                </div>
                <label htmlFor="login">Логин</label>
                <input
                    id="login"
                    type="text"
                    value={login}
                    onChange={(e) => setLogin(e.target.value)}
                />

                <label htmlFor="password">Пароль</label>
                <input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />

                <div className="display-flex align-items-center justify-content-end">
                    <Button type="button" className="text" onClick={toggleFormType}>{secondaryBtnTitle}</Button>
                    <Button type="submit">{actionBtnTitle}</Button>
                </div>
            </div>
        </form>
    );
}
