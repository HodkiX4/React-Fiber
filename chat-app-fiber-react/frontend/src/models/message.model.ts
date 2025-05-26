export interface IMessage {
  id: string;
  roomId: string;
  senderId: string;
  content: string;
  createdAt: string;
}

export interface ISendMessageRequest {
  roomId: string;
  content: string;
}

export interface IDeleteMessageResponse {
  message: string;
}
export interface ISendMessageResponse {
  message: IMessage
}

export interface IGetMessagesRequest {
  roomId: string
}

export interface IGetMessagesResponse {
  messages: IMessage[]
}