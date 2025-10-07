import { useEffect, useState, useSyncExternalStore } from "react";
import { SidebarItem } from "../SidebarItem/SidebarItem";
import Button from "../universal/Button/Button";
import type { TestCaseGroup } from "../../models/TestCase";
import { http } from "../../http";
import { appStore } from "../../Storage";
import { useLocation } from "react-router-dom";

export const Sidebar = () => {
    const [groupList, setGroupList] = useState<TestCaseGroup[]>([]);
    const store = useSyncExternalStore(appStore.subscribe, appStore.getSnapshot);
    const pathname = useLocation().pathname;
    const currentGroup = pathname.replace("/", "");

    const logout = async () => {
        await http.request("/auth/logout", {
            method: "POST"
        });

        http.deleteTokens();
        appStore.discard();
    }

    const saveFile = (filename: string, content: string) => {
        const binaryString = atob(content);

        const bytes = new Uint8Array(binaryString.length);
        for (let i = 0; i < binaryString.length; i++) {
            bytes[i] = binaryString.charCodeAt(i);
        }

        const blob = new Blob([bytes], { type: "application/octet-stream" });

        const link = document.createElement("a");
        link.href = URL.createObjectURL(blob);
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(link.href);
    }

    const getGroups = async () => {
        const response = await http.request<{ groups: TestCaseGroup[] }>("/groups/get");
        if (response.status === 200) {
            setGroupList(response.body.groups);
        }
    }

    const createGroup = async () => {
        const response = await http.request<TestCaseGroup>("/groups/add", {
            method: "POST"
        });

        if (response.status === 200) {
            setGroupList([...groupList, response.body]);
        }
    }

    const deleteGroup = async (uuid: string) => {
        const response = await http.request("/groups/delete", {
            method: "POST",
            body: { uuid }
        });

        if (response.status === 200) {
            setGroupList((g) => g.filter((el) => el.uuid !== uuid));
        }
    }

    const exportGroup = async (uuid: string, name: string, type: "excel" | "zephyr") => {
        let filename: string;
        let path: string;
        switch (type) {
            case "excel":
                path = "/groups/excel/" + uuid
                filename = name + ".xlsx"
                break;
            case "zephyr":
                path = "/groups/zephyr/" + uuid
                filename = name + ".json"
                break;
            default:
                return;
        }

        const response = await http.request<{ content: string }>(path);
        if (response.status === 200) {
            saveFile(filename, response.body.content);
        }
    }

    const renameGroup = async (name: string, uuid: string) => {
        const response = await http.request<TestCaseGroup>("/groups/rename", {
            method: "POST",
            body: { name, uuid }
        });

        if (response.status === 200) {
            setGroupList((g) => g.map((el) => el.uuid !== uuid ? el : response.body));
        }
    }

    useEffect(() => {
        getGroups();
    }, []);

    return (
        <div className="sidebar flex-wooden display-flex flex-direction-column overflow-y-hidden">
            <div className="case-group-list flex-rubber display-flex flex-direction-column overflow-y-auto">
                {groupList.map((g) => (
                    <SidebarItem
                        key={g.uuid}
                        isActive={currentGroup === g.uuid}
                        userUUID={store.userUUID}
                        uuid={g.uuid}
                        title={g.name}
                        creator={g.creator_uuid}
                        onDelete={deleteGroup}
                        onExport={exportGroup}
                        onRename={renameGroup}
                    />
                ))}
            </div>

            <div className="display-flex justify-content-end">
                <Button icon="logout" onClick={logout}>Выход</Button>
                <Button icon="create_new_folder" onClick={createGroup}>Создать новую группу</Button>
            </div>
        </div>
    );
};
