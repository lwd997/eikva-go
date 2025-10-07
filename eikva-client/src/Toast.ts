import { createElement } from "react";
import toast from "react-hot-toast";
import Icon from "./components/universal/Icon/Icon";

export const createErrorToast = (text: string) => {
    toast(text, {
        duration: 2500,
        className: "eikva-toast-error",
        icon: createElement(Icon, { name: "bug_report" })
    });
};

