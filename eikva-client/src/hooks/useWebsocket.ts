import { useEffect } from "react";
import { http } from "../http";
import { createErrorToast } from "../Toast";

export const useWebsocketUpdate = (onUpdate: (type: string, updateList: string[]) => void) => {
    useEffect(() => {
        const socket = new WebSocket(
            `${window.location.protocol === "http:" ? "ws" : "wss"}://${window.location.host}/ws`
        );

        socket.onopen = () => {
            socket.send(JSON.stringify({
                access_token: http.accessToken,
                type: 'auth'
            }));
        };

        socket.onmessage = (event) => {
            try {
                const mesasge: { type: string; uuid: string[] } = JSON.parse(event.data);
                if (
                    mesasge.type !== "test-case-update" &&
                    mesasge.type !== "upload-update"
                ) {
                    throw new Error("Неизвестный тип сообщения WebSocket");
                }

                onUpdate(mesasge.type, mesasge.uuid);
            } catch (error) {
                const errorMessage = error instanceof Error
                    ? error.message
                    : String(error)

                createErrorToast(errorMessage);
            }
        };

        socket.onclose = (event) => {
            if (event.code !== 1000) {
                createErrorToast(event.reason);
            }
        };

        socket.onerror = () => {
            createErrorToast("Ошибка при подключении к WebSocket");
        };

        return () => {
            socket.close(1000);
        }
    }, []);
};

