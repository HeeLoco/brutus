export interface IProvider {
    Id: string;
    Name: string;
    Entities: IEntity[];
}

export interface IEntity {
    Id: string;
    Name: string;
    Create: string;
    CreateMethod: string;
    GetConfig: string;
    GetConfigMethod: string;
}