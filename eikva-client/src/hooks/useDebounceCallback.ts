import { useRef, useEffect, useCallback } from "react";

export function useDebounceCallback<T extends (...args: any[]) => void>(
    fn: T,
    delay: number
) {
    const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const debounced = useCallback(
        (...args: Parameters<T>) => {
            if (timerRef.current) {
                clearTimeout(timerRef.current);
            }
            timerRef.current = setTimeout(() => fn(...args), delay);
        },
        [fn, delay]
    );

    useEffect(() => {
        return () => {
            if (timerRef.current) {
                clearTimeout(timerRef.current);
            }
        };
    }, []);

    return debounced;
}
