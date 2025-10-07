import { useEffect, useState } from "react";
import type { Step, TestCase } from "../../models/TestCase";
import { useDebounceCallback } from "../../hooks/useDebounceCallback";
import { http } from "../../http";
import { Steps } from "../../pages/Group/Steps/Steps";
import Badge from "../universal/Badge/Badge";
import Button from "../universal/Button/Button";
import { PreloadOverlay } from "../universal/PreloadOverlay/PreloadOverlay";
import { ErrorStatus, LoadingStatus } from "../../models/Status";

type CaseListItemProps = TestCase & {
    userUUID: string | null;
    onDelete: (tcUUID: string) => void;
}

export const CaseListItem = ({ userUUID, onDelete, ...testCase }: CaseListItemProps) => {
    const [state, setState] = useState<TestCase>(testCase);
    const [steps, setSteps] = useState<Step[]>([]);
    const [isStepsLoaded, setIsStepsLoaded] = useState(false);
    const [isDeleteConfirm, setIsDeleteConfirm] = useState(false);

    const [isExpanded, setIsExpanded] = useState(false);
    const isCreatedByUser = state.creator_uuid === userUUID;

    useEffect(() => {
        const listener = (e: CustomEvent<string[]>) => {
            if (e.detail.includes(state.uuid)) {
                getSelf();
            }
        }

        window.addEventListener('test-case-update', listener);
        return () => {
            window.removeEventListener('test-case-update', listener);
        }

    }, [state.uuid, isExpanded, isStepsLoaded]);

    const getSelf = async () => {
        if (isExpanded || isStepsLoaded) {
            console.log("??")
            getSteps();
        }

        const response = await http.request<TestCase>("/test-cases/get/" + state.uuid);
        if (response.status === 200) {
            setState(response.body);
        }
    };

    const badges = [
        { label: "Автор", value: state.creator },
        { label: "Дата создания", value: state.created_at }
    ];

    if (isCreatedByUser) {
        badges[0].value += " (Это вы)"
    }

    const toggleDeleteConfirm = () => setIsDeleteConfirm(!isDeleteConfirm);

    const toggleExpand = () => {
        if (!isStepsLoaded) {
            getSteps();
        }

        setIsExpanded(!isExpanded);
    }

    const onChange = <T extends keyof TestCase>(field: T, value: TestCase[T]) => {
        if (!isCreatedByUser) {
            return;
        }

        setState((s) => {
            s = { ...s, [field]: value };
            sync(s);
            return s;
        });
    }

    const sync = useDebounceCallback((tc: TestCase) => {
        http.request("/test-cases/update", {
            method: "POST",
            body: {
                description: tc.description,
                post_condition: tc.post_condition,
                pre_condition: tc.pre_condition,
                name: tc.name,
                uuid: tc.uuid,
                source_ref: tc.source_ref
            }
        });
    }, 500);

    const getSteps = async () => {
        const response = await http.request<{ steps: Step[] }>("/test-cases/get-steps/" + state.uuid);

        if (response.status === 200) {
            setSteps(response.body.steps);
            setIsStepsLoaded(true);
        }
    }

    const addTestCase = async () => {
        const response = await http.request<Step>("/steps/add", {
            method: "POST",
            body: {
                test_case: state.uuid
            }
        });

        if (response.status === 200) {
            setSteps([...steps, response.body]);
        }
    }

    let className = "card";
    if (state.status === ErrorStatus) {
        className += " error";
        state.name = "При генерации тест-кейса произошла ошибка"
    }

    return (
        <div className={className}>
            {state.status === LoadingStatus &&
                <PreloadOverlay/>
            }
            <div className="margin-bottom-1">
                <div className="card-top display-flex align-items-center justify-content-space-between">
                    <h1 className="card-heading"><span>{state.name}</span></h1>
                    <div className="card-top-actions display-flex align-items-center">
                        {isCreatedByUser &&
                            <>
                                {isDeleteConfirm
                                    ? (
                                        <>
                                            <Button onClick={() => onDelete(state.uuid)} icon="delete_forever">Да, удалить!</Button>
                                            <Button onClick={toggleDeleteConfirm}>Не удалять</Button>
                                        </>
                                    )
                                    : <Button icon="delete" onClick={toggleDeleteConfirm} />
                                }
                            </>
                        }

                        {isExpanded
                            ? <Button onClick={toggleExpand} icon="keyboard_arrow_up" />
                            : <Button onClick={toggleExpand} icon="keyboard_arrow_down" />
                        }
                    </div>
                </div>
                <div className="display-flex align-items-center">
                    {badges.map((badge, i) => <Badge key={i} label={badge.label} value={badge.value} />)}
                </div>
            </div>
            {isExpanded &&
                <div className="display-flex flex-direction-column">
                    <div>
                        <input
                            disabled={!isCreatedByUser}
                            onChange={(e) => onChange("name", e.target.value)}
                            value={state.name}
                        />
                    </div>
                    <div>
                        <div className="label">Источник</div>
                        <input
                            disabled={!isCreatedByUser}
                            onChange={(e) => onChange("source_ref", e.target.value)}
                            value={state.source_ref}
                        />
                    </div>

                    <div>
                        <div className="label">Описание</div>
                        <textarea
                            rows={10}
                            className="textarea"
                            disabled={!isCreatedByUser}
                            onChange={(e) => onChange("description", e.target.value)}
                            value={state.description}
                        />
                    </div>

                    <div className="display-flex width-100">
                        <div className="width-50">
                            <div className="label">Предусловие</div>
                            <textarea
                                rows={5}
                                className="textarea"
                                disabled={!isCreatedByUser}
                                onChange={(e) => onChange("pre_condition", e.target.value)}
                                value={state.pre_condition}
                            />
                        </div>

                        <div className="width-50">
                            <div className="label">Постусловие</div>
                            <textarea
                                rows={5}
                                className="textarea"
                                disabled={!isCreatedByUser}
                                onChange={(e) => onChange("post_condition", e.target.value)}
                                value={state.post_condition}
                            />
                        </div>
                    </div>

                    <div className="margin-top-1">
                        <Steps
                            rewriteSteps={setSteps}
                            isCreatedByUser={isCreatedByUser}
                            addTestCase={addTestCase}
                            isStepsLoaded={isStepsLoaded}
                            retryHandler={getSteps}
                            steps={steps}
                        />
                    </div>
                </div>
            }
        </div>
    );
};
