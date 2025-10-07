import { Tooltip } from "react-tooltip";
import type { UploadedFile } from "../../models/File";
import Icon from "../universal/Icon/Icon";
import { PreloadOverlay } from "../universal/PreloadOverlay/PreloadOverlay";
import { LoadingStatus } from "../../models/Status";

type FileElementProps = UploadedFile & {
    currentUserUUID: string;
    onPreview: (uuid: string) => void;
    onDelete: (uuid: string) => void;
    onCompress: (uuid: string) => void;
    onSelect: (f: UploadedFile) => void;
    isSeltected: boolean;
}

const tooltip = (content: string, id: string) => {
    return {
        "data-tooltip-content": content,
        "data-tooltip-id": id,
        "data-tooltip-placer": "top",
    }
}

export const FileElement = ({
    currentUserUUID,
    onPreview,
    onDelete,
    onCompress,
    onSelect,
    isSeltected,
    ...file
}: FileElementProps) => {
    let className = "file fc";
    if (isSeltected) {
        className += " selected";
    }

    return (
        <div className={className}>
            {file.status === LoadingStatus && <PreloadOverlay type="dots" /> }
            <div {...tooltip("Прикрепить", "add")} className="file" onClick={() => onSelect(file)}>
                <Icon name="docs" />
                <div>{file.name}</div>
            </div>
            <Icon {...tooltip("Уменьшить объем", "compress")} name="wand_stars" onClick={() => onCompress(file.uuid)} />
            <Icon {...tooltip("Просмотр содержимого", "preview")} name="search" onClick={() => onPreview(file.uuid)} />
            {file.creator === currentUserUUID &&
                <>
                    <Icon {...tooltip("Удалить", "delete")} name="close" onClick={() => onDelete(file.uuid)} />
                    <Tooltip id="delete" />
                </>
            }
            <Tooltip id="compress" />
            <Tooltip id="preview" />
            <Tooltip id="add" />
        </div>
    );
}

