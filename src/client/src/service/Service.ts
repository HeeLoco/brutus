import { IEntity } from "../interfaces/IProvider";

export async function fetchEndpoints() {
    const v = await  fetch('/api/providers', {
        method: 'GET'
    });
    if(!v.ok) {
        return v.json();
    }
}