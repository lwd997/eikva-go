import { Main } from "./pages/Main";
import { useSyncExternalStore } from "react";
import { Login } from "./pages/Login";
import { appStore } from "./Storage";
import { Toaster } from "react-hot-toast";

const App = () => {
    const store = useSyncExternalStore(appStore.subscribe, appStore.getSnapshot);
    return (
        <>
            <Toaster position="top-right" />
            {!store.isAuthorized ? <Login /> : <Main />}
        </>
    );
};

export default App;
