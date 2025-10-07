import type { Status } from "./Status";

export interface TestCaseGroup {
    id: number;
    uuid: string;
    name: string;
    status: Status;
    creator: string;
    creator_uuid: string;
}

export interface TestCase {
    id: number;
    uuid: string;
    created_at: string;
    creator: string;
    creator_uuid: string;
    test_case_group: string;
    status: Status;
    name: string;
    pre_condition: string;
    post_condition: string;
    description: string;
    source_ref: string;
}

export type TestCaseUpdatePayload = Pick<
    TestCase,
    "uuid" |
    "name" |
    "pre_condition" |
    "post_condition" |
    "description"
>

export interface Step {
    id: number;
    uuid: string;
    num: number;
    creator: string;
    creator_uuid: string;
    test_case: string;
    created_at: string;
    status: Status;
    data: string;
    description: string;
    expected_result: string;
}
