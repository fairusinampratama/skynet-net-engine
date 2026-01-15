import { twMerge } from "tailwind-merge";
import clsx from "clsx";

export function Skeleton({ className, ...props }) {
    return (
        <div
            className={twMerge(clsx("animate-pulse rounded-md bg-slate-200/80", className))}
            {...props}
        />
    );
}
