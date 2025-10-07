import type { Status } from "./Status";

export interface UploadedFile {
    uuid: string;
    name: string;
    content?: string;
    status: Status;
    token_count: number;
    creator: string;
    test_case_group: string;
}

