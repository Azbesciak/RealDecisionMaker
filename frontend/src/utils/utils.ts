import {ChangeEvent} from "react";
import {Collection} from "./ValuesContainerComponent";

export function handleInputValueChange<T extends string>(callback: (value: T) => void) {
    return (event: ChangeEvent<{ value: unknown }>) => callback(event.target.value as T);
}

export function remapCollection<T, R>(obj: Collection<T>, mapper: (k: string, v: T) => R): Collection<R> {
    return Object.fromEntries(Object.entries(obj).map(([k, v]) => [k, mapper(k, v)]))
}

export function camelCaseToNormal(value: string) {
    if (!isString(value)) return value;
    return value
        .replace(/([A-Z]+[^A-Z])/g, match => ` ${match}`)
        .replace(/([A-Z]+)([A-Z][a-z]+)/g, (_, ...low) => `${low[0]} ${low[1]}`)
        .replace(/([a-z]+)([\dA-Z]+)/g, (_, ...low) => `${low[0]} ${low[1]}`)
        .trim()
        .replace(/([A-Z][a-z]+)/g, match => match.toLowerCase())
        .replace(/^./, match => match.toUpperCase());
}

export function isString(value: any): value is string {
    return typeof value === "string";
}