export interface ItemValue<V> {
    value: V;
    onChange: (v: V) => void;
}

export interface Identifiable {
    id: string;
}