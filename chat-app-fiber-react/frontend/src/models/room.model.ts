export interface IRoom {
  id: string;
  name: string;
  createdAt: string;
  members: string[];
  isPrivate: boolean;
  admins: string[];
  moderators: string[];
}

export interface IGetRoomsResponse {
  rooms: IRoom[];
}

export interface ICreateRoomRequest {
  name: string;
  isPrivate: boolean;
  members?: string[];
}

export interface ICreateRoomResponse {
  message: string;
  room: IRoom;
}

export interface IDeleteRoomResponse {
  message: string;
}