import { Route, Routes } from "react-router-dom";
import { Sidebar } from "../components/Sidebar/Sidebar";
import { Group } from "./Group";
import { useEffect, useState } from "react";
import { appStore } from "../Storage";
import { http } from "../http";
import type { WhoAmIResponse } from "../models/Auth";
import { Home } from "./Home";
import { PreloadOverlay } from "../components/universal/PreloadOverlay/PreloadOverlay";

export const Main = () => {
    const [isChecking, setIsChecking] = useState(true);

    useEffect(() => {
        const checkSession = async () => {
            const response = await http.request<WhoAmIResponse>("/auth/whoami")
            const isOk = response.status === 200;
            appStore.updatePart({
                isAuthorized: isOk,
                userUUID: isOk ? response.body.uuid : null,
                userLogin: isOk ? response.body.uuid : null
            });

            if (isOk) {
                setIsChecking(false);
            }
        }

        checkSession();
    }, []);

    if (isChecking) {
        return (
            <div className="display-flex width-100 height-100">
                <PreloadOverlay />
            </div>
        )
    }

    return (
        <div className="display-flex width-100 height-100">
            <Sidebar />
            <div className="content display-flex flex-direction-column align-items-center">
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/:group" element={<Group />} />
                </Routes>
            </div>
        </div>
    );
};
