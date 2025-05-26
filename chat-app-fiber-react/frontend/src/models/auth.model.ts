export interface ILoginRequest {
  email: string;
  password: string;
}

export interface ILoginResponse {
  message: string;
  token: string;
}

export interface ISignupRequest {
  email: string;
  username: string;
  password: string;
  passwordAgain?: string;
}

export interface ISignupResponse {
  message: string;
}