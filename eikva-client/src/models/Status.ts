export const NoStatus = 'none';
export const LoadingStatus = 'loading';
export const ErrorStatus = 'error';

export type Status = (
    typeof NoStatus |
    typeof LoadingStatus |
    typeof ErrorStatus
);
