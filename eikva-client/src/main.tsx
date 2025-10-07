import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import.meta.glob('./styles/*.css', { eager: true });
import App from "./App.tsx";

declare global {
    interface WindowEventMap {
        "test-case-update": CustomEvent<string[]>;
        "upload-update": CustomEvent<string[]>;
    }
}

const root = document.getElementById("root")
if (!root) {
    throw new Error("На странице отсуствует #root");
}

root.className = "width-100 height-100";

createRoot(root).render(
    <BrowserRouter>
        <App />
    </BrowserRouter>
);
