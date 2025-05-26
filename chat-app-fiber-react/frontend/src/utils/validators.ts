import { ILoginRequest, ISignupRequest } from "../models/auth.model";

export const isValidUsername = (username: string): boolean => {
  const trimmed = username.trim();
  const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/;
  return usernameRegex.test(trimmed);
};

export const isValidEmail = (email: string): boolean => {
  const trimmed = email.trim().toLowerCase();
  const emailRegex = /^[a-zA-Z0-9._%+/-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(trimmed);
};

export const isValidPassword = (password: string): boolean => {
  const trimmed = password.trim();
  const passwordRegex = /^[A-Za-z0-9!@#$%^&*()_+={}[\]:;"'<>,.?/-]{8,}$/;
  if (!passwordRegex.test(trimmed)) return false;

  const hasLetter = /[A-Za-z]/.test(trimmed);
  const hasDigit = /\d/.test(trimmed);
  const hasSpecial = /[!@#$%^&*()_+={}[\]:;"'<>,.?/-]/.test(trimmed);

  return hasLetter && hasDigit && hasSpecial;
};

export const validateSignupInput = (input: ISignupRequest): string | null => {
  if (!isValidUsername(input.username)) {
    return "Username must be alphanumeric and between 3 and 20 characters long";
  }

  if (!isValidEmail(input.email)) {
    return "Invalid email format";
  }

  if (!isValidPassword(input.password)) {
    return "Password must be at least 8 characters long and contain at least one letter, one number, and one special character";
  }

  if(input.password !== input.passwordAgain) {
    return "Passwords don't match"
  }
  
  return null;
};

export const validateLoginInput = (input: ILoginRequest): string | null => {
  if(
    input.email.trim().length == 0 ||
    input.password.trim().length == 0
  ) {
    return "Please fill out all fields"
  }
  return null
}