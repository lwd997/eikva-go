import { useEffect, useState, useSyncExternalStore } from "react";
import { useParams } from "react-router-dom"
import type { TestCase } from "../../models/TestCase";
import { CaseListItem } from "../../components/CaseListItem/CaseListItem";
import { http } from "../../http";
import Button from "../../components/universal/Button/Button";
import { appStore } from "../../Storage";
import { InputArea } from "../../components/InputArea";
import { useWebsocketUpdate } from "../../hooks/useWebsocket";
import "./Group.css";

export const Group = () => {
    const { group } = useParams();
    const { userUUID } = useSyncExternalStore(appStore.subscribe, appStore.getSnapshot);
    const [testCaseList, setTestCaseList] = useState<TestCase[]>([]);

    useWebsocketUpdate((type, updateList) => {
        window.dispatchEvent(new CustomEvent(type, {
            detail: updateList
        }));
    });

    const getTestCases = async () => {
        if (!group) {
            return;
        }

        const response = await http.request<{ test_cases: TestCase[] }>("/groups/get-test-cases/" + group);

        if (response.status === 200) {
            setTestCaseList(response.body.test_cases);
        }
    }

    const addTestCase = async () => {
        if (!group) {
            return;
        }

        const response = await http.request<TestCase>("/test-cases/add", {
            method: "POST",
            body: {
                test_case_group: group
            }
        });

        if (response.status === 200) {
            setTestCaseList([...testCaseList, response.body]);
        }
    }

    const deleteTestCase = async (testCaseUUID: string) => {
        const response = await http.request<TestCase>("/test-cases/delete", {
            method: "POST",
            body: {
                uuid: testCaseUUID
            }
        });

        if (response.status === 200) {
            setTestCaseList(testCaseList.filter((tc) => tc.uuid !== testCaseUUID));
        }
    }

    const addManyTestCases = (tcList: TestCase[]) => {
        setTestCaseList([...testCaseList, ...tcList]);
    }

    useEffect(() => {
        getTestCases();
    }, [group]);

    if (!group) {
        return null;
    }

    return (
        <>
            <div className="group display-flex flex-direction-column">
                {testCaseList.map((tc) => (
                    <CaseListItem
                        key={tc.uuid}
                        userUUID={userUUID}
                        onDelete={deleteTestCase}
                        {...tc}
                    />
                ))}
            </div>
            <InputArea
                currentUserUUID={userUUID!}
                addTestCases={addManyTestCases}
                testCaseGroupUUID={group}
            />
            <Button className="add-test-case" icon="add" onClick={addTestCase}>Добавить новый тест-кейс</Button>
        </>
    );
}
