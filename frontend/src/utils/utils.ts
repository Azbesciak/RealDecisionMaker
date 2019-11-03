import {ChangeEvent} from "react";
import {Collection} from "./ValuesContainerComponent";

export function handleInputValueChange<T extends string>(callback: (value: T) => void) {
    return (event: ChangeEvent<{ value: unknown }>) => callback(event.target.value as T);
}

export function remapCollection<T, R>(obj: Collection<T>, mapper: (k: string, v: T) => R): Collection<R> {
    return Object.fromEntries(Object.entries(obj).map(([k, v]) => [k, mapper(k, v)]))
}