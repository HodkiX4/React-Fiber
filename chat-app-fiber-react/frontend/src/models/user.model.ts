// types/User.ts

export type UserStatus = 'online' | 'offline' | 'busy' | 'away'; // Adjust as needed

export interface IUser {
  id: string;
  username: string;
  email: string;
  avatarUrl?: string;
  status?: UserStatus;
  friends?: string[];
  createdAt?: string;
}

export interface IGetUsersResponse {
  users: IUser[]
}