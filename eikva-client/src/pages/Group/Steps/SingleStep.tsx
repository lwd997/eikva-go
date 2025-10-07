import { useState } from "react";
import type { Step } from "../../../models/TestCase";
import { http } from "../../../http";
import { useDebounceCallback } from "../../../hooks/useDebounceCallback";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import Button from "../../../components/universal/Button/Button";
import Icon from "../../../components/universal/Icon/Icon";

type SingleStepProps = Step & {
    isCreatedByUser: boolean;
    onDelete: (uuid: string) => void;
}

export const SingleStep = ({ isCreatedByUser, onDelete, ...step }: SingleStepProps) => {
    const [state, setState] = useState(step);
    const [isDeleteConfirm, setIsDeleteConfirm] = useState(false);

    const {
        attributes,
        listeners,
        setNodeRef,
        transform,
        transition,
    } = useSortable({ id: step.uuid });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
    };

    const toggleDeleteConfirm = () => setIsDeleteConfirm(!isDeleteConfirm);

    const onChange = <T extends keyof Step>(field: T, value: Step[T]) => {
        if (!isCreatedByUser) {
            return;
        }

        setState((s) => {
            s = { ...s, [field]: value };
            sync(s);
            return s;
        });
    }

    const sync = useDebounceCallback((s: Step) => {
        http.request("/steps/update", {
            method: "POST",
            body: {
                uuid: s.uuid,
                data: s.data,
                expected_result: s.expected_result,
                description: s.description
            }
        });
    }, 500);

    return (
        <div
            ref={setNodeRef}
            style={style}
            className="single-step"
        >
            <div
                {...attributes}
                {...listeners} className="grabber">
                <Icon name="drag_indicator"/>
            </div>
            <div className="step-content display-flex flex-direction-column">
                <div className="display-flex align-items-center justify-content-space-between">
                    <h2>No: {step.num}</h2>

                    <div className="display-flex align-items-center justify-content-space-between">
                        {/*badges.map((badge, i) => <Badge key={i} label={badge.label} value={badge.value} />)*/}
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
                    </div>
                </div>
                <div className="display-flex width-100">
                    <div className="width-50">
                        <div className="label">Описание</div>
                        <textarea
                            rows={6}
                            className="textarea"
                            disabled={!isCreatedByUser}
                            onChange={(e) => onChange("description", e.target.value)}
                            value={state.description}
                        />
                    </div>

                    <div className="width-50">
                        <div className="label">Ожидаемый результат</div>
                        <textarea
                            rows={6}
                            className="textarea"
                            disabled={!isCreatedByUser}
                            onChange={(e) => onChange("expected_result", e.target.value)}
                            value={state.expected_result}
                        />
                    </div>
                </div>


                <div>
                    <div className="label">Данные</div>
                    <textarea
                        rows={6}
                        className="textarea"
                        disabled={!isCreatedByUser}
                        onChange={(e) => onChange("data", e.target.value)}
                        value={state.data}
                    />
                </div>
            </div>
        </div>
    )
}
