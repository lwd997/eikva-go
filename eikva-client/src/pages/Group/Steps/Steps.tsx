import { useState } from "react"
import type { Step } from "../../../models/TestCase";
import Button from "../../../components/universal/Button/Button";
import { SingleStep } from "./SingleStep";
import { closestCenter, DndContext, KeyboardSensor, PointerSensor, useSensor, useSensors, type DragEndEvent } from "@dnd-kit/core";
import { arrayMove, SortableContext, sortableKeyboardCoordinates, verticalListSortingStrategy } from "@dnd-kit/sortable";
import { http } from "../../../http";

import "./Steps.css";

interface StepsProps {
    steps: Step[];
    isStepsLoaded: boolean;
    retryHandler: () => void;
    addTestCase: () => void;
    isCreatedByUser: boolean;
    rewriteSteps: (s: Step[]) => void;
}

export const Steps = ({
    steps,
    isStepsLoaded,
    isCreatedByUser,
    retryHandler,
    addTestCase,
    rewriteSteps
}: StepsProps) => {
    const [isExpanded, setIsExpanded] = useState(true);
    const sensors = useSensors(
        useSensor(PointerSensor),
        useSensor(KeyboardSensor, {
            coordinateGetter: sortableKeyboardCoordinates,
        })
    );

    const toggleExpaned = () => {
        setIsExpanded(!isExpanded);
    }

    const swapSteps = async (first: string, second: string) => {
        await http.request("/steps/swap", {
            method: "POST",
            body: {
                first,
                second
            }
        });

    };

    const deleteStep = async (stepUUID: string) => {
        const toDelete = steps.findIndex((s) => s.uuid === stepUUID);
        if (toDelete === -1) {
            return;
        }

        const response = await http.request("/steps/delete", {
            method: "POST",
            body: { uuid: stepUUID }
        });

        if (response.status === 200) {
            retryHandler();
        }
    }

    const handleDragEnd = async (event: DragEndEvent) => {
        const { active, over } = event;

        if (active.id !== over?.id) {
            const firstIndex = steps.findIndex((s) => s.uuid === active.id);
            const secondIndex = steps.findIndex((s) => s.uuid === over!.id);
            if (firstIndex !== -1 && secondIndex !== -1) {
                const firstNum = steps[firstIndex].num;
                const secondNum = steps[secondIndex].num;
                steps[firstIndex].num = secondNum;
                steps[secondIndex].num = firstNum;
                rewriteSteps(arrayMove(steps, firstIndex, secondIndex));
                await swapSteps(steps[firstIndex].uuid, steps[secondIndex].uuid);
                retryHandler()
            }
        }
    }

    return (
        <div className="display-flex flex-direction-column">
            <div>
                <Button
                    className="text row-reverse"
                    onClick={toggleExpaned}
                    icon={isExpanded ? "keyboard_arrow_up" : "keyboard_arrow_down"}
                >
                    Шаги
                </Button>

            </div>

            {isExpanded &&
                <div className="steps display-flex flex-direction-column">
                    {!isStepsLoaded &&
                        <div>
                            Не удалось загрузить шаги
                            <Button
                                onClick={retryHandler}
                                icon="footprint"
                            >
                                Повторить попытку
                            </Button>

                        </div>
                    }

                    {isStepsLoaded &&
                        <>
                            {steps.length
                                ? (
                                    <>
                                        {isCreatedByUser
                                            ? (
                                                <DndContext
                                                    sensors={sensors}
                                                    collisionDetection={closestCenter}
                                                    onDragEnd={handleDragEnd}
                                                >
                                                    <SortableContext
                                                        items={steps.map((s) => s.uuid)}
                                                        strategy={verticalListSortingStrategy}
                                                    >
                                                        {steps.map((s,) => (
                                                            <SingleStep
                                                                key={s.uuid}
                                                                isCreatedByUser={isCreatedByUser}
                                                                onDelete={deleteStep}
                                                                {...s}
                                                            />
                                                        ))}
                                                    </SortableContext>
                                                </DndContext>
                                            )
                                            : (
                                                steps.map((s,) => (
                                                    <SingleStep
                                                        key={s.uuid}
                                                        isCreatedByUser={isCreatedByUser}
                                                        onDelete={deleteStep}
                                                        {...s}
                                                    />
                                                ))
                                            )
                                        }
                                    </>
                                )
                                : "Пока не добавлено ни одного шага"
                            }
                        </>
                    }

                    {isCreatedByUser &&
                        <div className="display-flex justify-content-end">
                            <Button className="text" onClick={addTestCase} icon="add">Добавить шаг</Button>
                        </div>
                    }
                </div>
            }
        </div>
    );
}

